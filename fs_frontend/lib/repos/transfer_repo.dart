import 'dart:async';
import 'dart:convert';
import 'dart:io';

import 'package:fs_frontend/constants/api_constants.dart';
import 'package:fs_frontend/exceptions/exceptions.dart';
import 'package:fs_frontend/models/transfer_task.dart';
import 'package:fs_frontend/repos/api_repo.dart';
import 'package:http/http.dart' as http;

class TransferRepo extends ApiRepo {
  TransferRepo({required super.secureStorage}) {
    _listener();
  }

  void _listener() {
    _controller.stream.listen((event) {
      if (event is UploadTask) {
        _upload(event);
      } else if (event is DownloadTask) {
        _download(event);
      }
    });
  }

  Future<void> _upload(UploadTask task) async {
    try {
      final token = await secureStorage.read(key: 'token');
      if (token == null || token.isEmpty) {
        throw UnauthorizedError();
      }

      // Initial upload request
      String uploadLink = "";
      final uri = Uri.parse(ApiConstants.fileQuery)
          .replace(queryParameters: {'location': task.parentLocation});
      final uploadRequestRes = await super.authRequest(
          requestType: RequestType.post,
          uri: uri,
          body: jsonEncode({'name': task.fileName, 'size': task.size}));
      if (uploadRequestRes.statusCode == 200) {
        uploadLink = jsonDecode(uploadRequestRes.body)['uploadLink'];
      } else if (uploadRequestRes.statusCode == 403) {
        throw UnauthorizedError();
      } else {
        throw ServerError();
      }

      // File upload
      final client = http.Client();
      http.Request uploadRequest;

      const chunkSize = 1024 * 1024;
      int startByte = 0;
      int endByte = chunkSize;
      int totalChunkNumber = (task.size / chunkSize).ceil();
      int chunkNumber = 1;

      while (startByte < task.size) {
        Stream<List<int>> chunkStream;
        if (endByte >= task.size) {
          // The last chunk might be smaller than the chunk size, adjust the length accordingly.
          chunkStream = task.file.openRead(startByte);
          endByte = task.size;
        } else {
          chunkStream = task.file.openRead(startByte, endByte);
        }

        final chunkData =
            await chunkStream.fold<List<int>>([], (previous, element) {
          previous.addAll(element);
          return previous;
        });

        uploadRequest = http.Request(
            'POST', Uri.parse('${ApiConstants.upload}/$uploadLink'));
        uploadRequest.headers['Content-Type'] = 'application/octet-stream';
        uploadRequest.headers['Chunk-Total'] = totalChunkNumber.toString();
        uploadRequest.headers['Chunk-Current'] = chunkNumber.toString();
        uploadRequest.headers['Authorization'] = 'bearer $token';
        uploadRequest.body = jsonEncode({'data': base64Encode(chunkData)});

        final uploadResponse = await client.send(uploadRequest);

        if (uploadResponse.statusCode == 200) {
          task.progressStreamController.sink.add(endByte / task.size);
          task.progress = endByte / task.size;
        } else {
          throw Error();
        }

        startByte = endByte;
        endByte = startByte + chunkSize;
        chunkNumber++;
      }
      task.hasCompleted = true;
      client.close();
    } catch (_) {
      task.hasFailed = true;
    }
  }

  Future<void> _download(DownloadTask task) async {
    try {
      final token = await secureStorage.read(key: 'token');
      if (token == null || token.isEmpty) {
        throw UnauthorizedError();
      }
      // Initial download request
      String downloadLink = "";
      int chunkTotal = 0;
      final uri = Uri.parse(ApiConstants.fileQuery)
          .replace(queryParameters: {'location': task.fileInCloud.location});
      final downloadRequestRes =
          await super.authRequest(requestType: RequestType.get, uri: uri);
      if (downloadRequestRes.statusCode == 200) {
        final parsed = jsonDecode(downloadRequestRes.body);
        downloadLink = parsed['downloadLink'];
        chunkTotal = parsed['chunkTotal'];
        if (downloadLink.isEmpty) {
          throw ServerError();
        }
      } else if (downloadRequestRes.statusCode == 403) {
        throw UnauthorizedError();
      } else {
        throw ServerError();
      }

      // File download
      int chunkNumber = 1;
      final File file = File(task.downloadPathString);

      for (int i = 1; i <= chunkTotal; i++) {
        final downloadResponse = await http.get(
            Uri.parse('${ApiConstants.download}/$downloadLink'),
            headers: {
              'Authorization': 'bearer $token',
              'Chunk-Current': chunkNumber.toString(),
              'Chunk-Total': chunkTotal.toString(),
            });
        if (downloadResponse.statusCode == 200) {
          await file.writeAsBytes(downloadResponse.bodyBytes,
              mode: FileMode.append);
          task.progressStreamController.sink.add(chunkNumber / chunkTotal);
          task.progress = chunkNumber / chunkTotal;
        } else {
          throw Error();
        }
        chunkNumber++;
      }

      task.hasCompleted = true;
    } catch (_) {
      task.hasFailed = true;
    }
  }

  final _controller = StreamController<TransferTask>();

  // Stream<TransferTask> get stream => _controller.stream;
  Sink<TransferTask> get sink => _controller.sink;
}
