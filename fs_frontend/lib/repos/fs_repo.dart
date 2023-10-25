import 'dart:async';
import 'dart:convert';

import 'package:fs_frontend/constants/api_constants.dart';
import 'package:fs_frontend/exceptions/exceptions.dart';
import 'package:fs_frontend/models/directory.dart';
import 'package:fs_frontend/models/file.dart';
import 'package:fs_frontend/models/role.dart';
import 'package:fs_frontend/repos/api_repo.dart';

class FSRepo extends ApiRepo {
  FSRepo({required super.secureStorage});

  Future<DirectoryModel> getDirectory(String location) async {
    final uri = Uri.parse(ApiConstants.dirQuery)
        .replace(queryParameters: {'location': location});
    final res = await super.authRequest(requestType: RequestType.get, uri: uri);
    if (res.statusCode == 200) {
      final parsed = jsonDecode(res.body);
      return DirectoryModel.fromJson(parsed['directoryAndContents']);
    } else if (res.statusCode == 403) {
      final parsed = jsonDecode(res.body);
      String errCode = parsed['error']['code'];
      // String errDescription = parsed['error']['description'];
      if (errCode == 'permission-denied') {
        throw PermissionDenied();
      } else {
        throw ServerError();
      }
    } else {
      throw ServerError();
    }
  }

  Future<DirectoryModel> createDirectory(
      String newDirName, String location) async {
    final uri = Uri.parse(ApiConstants.dirQuery)
        .replace(queryParameters: {'location': location});
    final res = await super.authRequest(
        requestType: RequestType.put,
        uri: uri,
        body: jsonEncode({'newDirectoryName': newDirName}));
    if (res.statusCode == 201) {
      final parsed = jsonDecode(res.body);
      return DirectoryModel.fromJson(parsed['newDirectory']);
    } else if (res.statusCode == 400) {
      final parsed = jsonDecode(res.body);
      final errCode = parsed['error']['code'];
      final errDescription = parsed['error']['description'];
      if (errCode == 'dir-already-exists') {
        throw ResourceAlreadyExists(errDescription);
      } else {
        throw ServerError();
      }
    } else if (res.statusCode == 403) {
      throw UnauthorizedError();
    } else {
      throw ServerError();
    }
  }

  FutureOr<void> deleteDirectory(String location) async {
    final uri = Uri.parse(ApiConstants.dirQuery)
        .replace(queryParameters: {'location': location});
    final res =
        await super.authRequest(requestType: RequestType.delete, uri: uri);
    if (res.statusCode == 200) {
      return;
    } else if (res.statusCode == 400) {
      final parsed = jsonDecode(res.body);
      final errCode = parsed['error']['code'];
      final errDescription = parsed['error']['description'];
      if (errCode == 'dir-not-empty') {
        throw ChildrenExists(description: errDescription);
      } else {
        throw ServerError();
      }
    } else if (res.statusCode == 403) {
      throw UnauthorizedError();
    } else {
      throw ServerError();
    }
  }

  FutureOr<void> deleteFile(String location) async {
    final uri = Uri.parse(ApiConstants.fileQuery)
        .replace(queryParameters: {'location': location});
    final res =
        await super.authRequest(requestType: RequestType.delete, uri: uri);
    if (res.statusCode == 200) {
      return;
    } else if (res.statusCode == 403) {
      throw UnauthorizedError();
    } else {
      throw ServerError();
    }
  }

  FutureOr<List<dynamic>> getShared(workspaceName) async {
    final uri = Uri.parse(ApiConstants.sharedContent)
        .replace(queryParameters: {'workspace': workspaceName});
    final res = await authRequest(requestType: RequestType.get, uri: uri);
    if (res.statusCode == 200) {
      final parsed = jsonDecode(res.body);
      final jsonList = parsed['sharedContent'] as List;
      final contents = [];
      for (var element in jsonList) {
        if (element['type'] == 'file') {
          contents.add(FileModel.fromJson(element));
        } else if (element['type'] == 'directory') {
          contents.add(DirectoryModel.fromJson(element));
        }
      }
      return contents;
    } else {
      throw ServerError();
    }
  }

  FutureOr<void> addRoleToFS(String location, Role role) async {
    final res = await authRequest(
        requestType: RequestType.post,
        uri: Uri.parse(ApiConstants.assignRoleFS),
        body: jsonEncode({
          'location': location,
          'roleId': role.id,
        }));
    if (res.statusCode == 201) {
      return;
    } else {
      throw ServerError();
    }
  }

  FutureOr<void> removeRoleFS(String location, Role role) async {
    final res = await authRequest(
        requestType: RequestType.delete,
        uri: Uri.parse(ApiConstants.assignRoleFS),
        body: jsonEncode({'roleId': role.id, 'location': location}));
    if (res.statusCode == 200) {
      return;
    } else {
      throw ServerError();
    }
  }

  Future<DirectoryModel> getDirectoryDetails(String location) async {
    final uri = Uri.parse(ApiConstants.dirDetails)
        .replace(queryParameters: {'location': location});
    final res = await authRequest(requestType: RequestType.get, uri: uri);
    if (res.statusCode == 200) {
      final parsed = jsonDecode(res.body);
      final dir = DirectoryModel.fromJson(parsed['directory']);
      final rolesParsed = parsed['roles'] as List<dynamic>;
      for (var element in rolesParsed) {
        dir.roles.add(Role.fromJson(element));
      }
      return dir;
    } else {
      throw ServerError();
    }
  }

  Future<List<Role>> getAllRoles(String workspaceName) async {
    final uri = Uri.parse(ApiConstants.allRolesQuery)
        .replace(queryParameters: {'workspaceName': workspaceName});
    final res = await super.authRequest(
        requestType: RequestType.get,
        uri: uri,
        body: jsonEncode({'workspaceName': workspaceName}));
    if (res.statusCode == 200) {
      final parsed = jsonDecode(res.body);
      List<Role> roles = [];
      final rolesJson = parsed['roles'] as List;
      for (var element in rolesJson) {
        roles.add(Role.fromJson(element));
      }
      return roles;
    } else {
      throw ServerError();
    }
  }
}
