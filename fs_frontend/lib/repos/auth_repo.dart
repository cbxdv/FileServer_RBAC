import 'dart:convert';

import 'package:fs_frontend/constants/api_constants.dart';
import 'package:fs_frontend/exceptions/exceptions.dart';
import 'package:fs_frontend/models/account.dart';
import 'package:fs_frontend/repos/api_repo.dart';
import 'package:http/http.dart' as http;

import '../models/workspace.dart';

class AuthRepo extends ApiRepo {
  AuthRepo({ required super.secureStorage });

  Future<Account?> checkAuth() async {
    final token = await secureStorage.read(key: "token");
    if (token == null || token.isEmpty) {
      return null;
    }
    final res = await http.get(
      Uri.parse(ApiConstants.checkAuth),
      headers: {'Authorization': 'bearer $token'}
    );
    if (res.statusCode == 200) {
      final parsed = jsonDecode(res.body);
      if (parsed['account']['isOwner'] == true) {
        return OwnerAccount.fromJson(parsed['account']);
      } else {
        final parsedAccount = parsed['account'];
        final usernameSplit = parsedAccount['username'].toString().split('@');
        final workspaceName = usernameSplit.last;
        final username = usernameSplit.first;
        return ServiceAccount(
            workspace: Workspace(name: workspaceName),
            id: parsedAccount['id'],
            name: parsedAccount['name'],
            username: username,
        );
      }
    }
    return null;
  }

  Future<OwnerAccount?> ownerAccountLogin({
    required String email,
    required String password,
  }) async {
    final res = await http.post(
      Uri.parse(ApiConstants.loginOwnerAccount),
      body: jsonEncode(<String, String>{'email': email, 'password': password}),
    );
    if (res.statusCode == 200) {
      final parsed = jsonDecode(res.body);
      await secureStorage.write(key: 'token', value: parsed['token']);
      return OwnerAccount.fromJson(parsed['account']);
    } else if (res.statusCode == 401) {
      final parsed = jsonDecode(res.body);
      String errCode = parsed['error']['code'];
      String errDescription = parsed['error']['description'];
      if (errCode == 'invalid-credentials') {
        throw InvalidCredentials(errDescription);
      } else {
        throw ServerError();
      }
    }
    return null;
  }

  Future<ServiceAccount?> serviceAccountLogin({
    required String username,
    required String password,
  }) async {
    final res = await http.post(
      Uri.parse(ApiConstants.loginServiceAccount),
      body: jsonEncode({ 'username': username, 'password': password })
    );
    if (res.statusCode == 200) {
      final parsed = jsonDecode(res.body);
      await secureStorage.write(key: 'token', value: parsed['token']);
      return ServiceAccount.fromJson(parsed['account']);
    } else if (res.statusCode == 401) {
      final parsed = jsonDecode(res.body);
      String errCode = parsed['error']['code'];
      String errDescription = parsed['error']['description'];
      if (errCode == 'invalid-credentials') {
        throw InvalidCredentials(errDescription);
      } else {
        throw ServerError();
      }
    }
    return null;
  }

  Future<OwnerAccount?> registerAccount({
    required String name,
    required String email,
    required String password,
  }) async {
    final res = await http.post(
      Uri.parse(ApiConstants.registerOwnerAccount),
      body: jsonEncode(<String, String>{"name": name, "email": email, "password": password}),
    );
    if (res.statusCode == 201) {
      return ownerAccountLogin(email: email, password: password);
    } else if (res.statusCode == 400) {
      final parsed = jsonDecode(res.body);
      String errCode = parsed['error']['code'];
      String errDescription = parsed['error']['description'];
      if (errCode == 'weak-password') {
        throw WeakPassword(errDescription);
      } else if (errCode == 'oa-already-exists') {
        throw AccountAlreadyExists(errDescription);
      }
    } else {
      throw ServerError();
    }
    return null;
  }

  Future<void> logout() async {
    secureStorage.deleteAll();
  }

  Future<void> changePassword({
    required String oldPassword, required String newPassword,
}) async {
    final res = await super.authRequest(
        requestType: RequestType.patch,
        uri: Uri.parse(ApiConstants.changePassword),
        body: jsonEncode({'oldPassword': oldPassword, 'newPassword': newPassword})
    );
    if (res.statusCode == 200) {
      return;
    } else {
      throw ServerError();
    }
  }
}
