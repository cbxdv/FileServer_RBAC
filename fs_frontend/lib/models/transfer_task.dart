import 'dart:async';
import 'dart:io';

import 'package:fs_frontend/models/file.dart';

abstract class TransferTask {
  late String fileName;
  int size;
  late StreamController<double> progressStreamController;
  late Stream<double> progressStream;
  late double progress;
  late bool hasCompleted;
  late bool hasFailed;

  TransferTask({ required this.size }) {
    progressStreamController = StreamController<double>();
    progressStream = progressStreamController.stream.asBroadcastStream();
    progress = 0;
    hasCompleted = false;
    hasFailed = false;
  }
}

class UploadTask extends TransferTask {
  final String parentLocation;
  final File file;

  UploadTask({required super.size, required this.parentLocation, required this.file}) {
    List<String> nameSplits = file.path.split('/');
    if (nameSplits.length == 1) {
      nameSplits = file.path.split('\\');
    }
    fileName = nameSplits.last;
  }
}

class DownloadTask extends TransferTask {
  final FileModel fileInCloud;
  String downloadPathString;
  DownloadTask({required this.fileInCloud, required this.downloadPathString, required super.size, required String fileName}) {
    this.fileName = fileName;
  }
}