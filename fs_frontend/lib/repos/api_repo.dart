import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:fs_frontend/exceptions/exceptions.dart';
import 'package:http/http.dart' as http;

enum RequestType {
  get, post, put, delete, patch
}

abstract class ApiRepo {
  final FlutterSecureStorage secureStorage;

  ApiRepo({required this.secureStorage});

  Future<http.Response> authRequest({required RequestType requestType, required Uri uri, String? body}) async {
    final token = await secureStorage.read(key: 'token');
    if (token == null || token.isEmpty) {
      throw UnauthorizedError();
    }
    switch(requestType) {
      case RequestType.get:
        return http.get(
          uri,
          headers: {'Authorization': 'bearer $token'},
        );
      case RequestType.post:
        return http.post(
          uri,
          headers: {'Authorization': 'bearer $token'},
          body: body,
        );
      case RequestType.put:
        return http.put(
          uri,
          headers: {'Authorization': 'bearer $token'},
          body: body,
        );
      case RequestType.delete:
        return http.delete(
          uri,
          headers: {'Authorization': 'bearer $token'},
          body: body,
        );
      case RequestType.patch:
        return http.patch(
          uri,
          headers: {'Authorization': 'bearer $token'},
          body: body,
        );
    }
  }
}