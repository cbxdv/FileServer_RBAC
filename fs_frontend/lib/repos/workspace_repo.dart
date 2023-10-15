import 'dart:convert';

import 'package:fs_frontend/constants/api_constants.dart';
import 'package:fs_frontend/exceptions/exceptions.dart';
import 'package:fs_frontend/models/account.dart';
import 'package:fs_frontend/models/role.dart';
import 'package:fs_frontend/models/workspace.dart';
import 'package:fs_frontend/repos/api_repo.dart';

class WorkspaceRepo extends ApiRepo {
  WorkspaceRepo({required super.secureStorage});

  Future<List<Workspace>> getAllWorkspaces() async {
    final res = await super.authRequest(
        requestType: RequestType.get,
        uri: Uri.parse(ApiConstants.workspacesOperations));
    if (res.statusCode == 200) {
      final parsed = jsonDecode(res.body);
      List<Workspace> workspaces = [];
      final workspacesJson = parsed['workspaces'] as List;
      for (var element in workspacesJson) {
        workspaces.add(Workspace.fromJson(element));
      }
      return workspaces;
    } else {
      throw ServerError();
    }
  }

  Future<Workspace> createNewWorkspace(String workspaceName) async {
    final token = await secureStorage.read(key: "token");
    if (token == null || token.isEmpty) {
      throw UnauthorizedError();
    }
    final res = await super.authRequest(
      requestType: RequestType.put,
      uri: Uri.parse(ApiConstants.workspacesOperations),
      body: jsonEncode({"workspaceName": workspaceName}),
    );
    if (res.statusCode == 201) {
      final parsed = jsonDecode(res.body);
      return Workspace.fromJson(parsed['workspace']);
    } else {
      if (res.statusCode == 400) {
        final parsed = jsonDecode(res.body);
        String errCode = parsed['error']['code'];
        if (errCode == 'workspace-exists') {
          throw WorkspaceAlreadyExists();
        }
      }
      throw ServerError();
    }
  }

  Future<ServiceAccount> createServiceAccount(Workspace workspace, String name,
      String username, String password) async {
    final res = await super.authRequest(
        requestType: RequestType.put,
        uri: Uri.parse(ApiConstants.workspaceAccountOperations),
        body: jsonEncode({
          'workspaceName': workspace.name,
          'name': name,
          'username': username,
          'password': password,
        }));
    if (res.statusCode == 201) {
      final parsed = jsonDecode(res.body);
      final parsedAcc = parsed['newServiceAccount'];
      return ServiceAccount(
          workspace: workspace,
          id: parsedAcc['id'],
          name: parsedAcc['name'],
          username: parsedAcc['username']
      );
    } else {
      throw ServerError();
    }
  }

  Future<void> updateWorkspaceName(String oldName, String newName) async {
    final res = await super.authRequest(
        requestType: RequestType.patch,
        uri: Uri.parse(ApiConstants.workspacesOperations),
        body: jsonEncode({
          'oldWorkspaceName': oldName,
          'newWorkspaceName': newName,
        }));
    if (res.statusCode == 200) {
      return;
    } else if (res.statusCode == 400) {
      final parsed = jsonDecode(res.body);
      String errCode = parsed['error']['code'];
      if (errCode == 'workspace-exists') {
        throw WorkspaceAlreadyExists();
      }
    } else {
      throw ServerError();
    }
  }

  Future<void> deleteWorkspace(String workspaceName) async {
    final res = await super.authRequest(
        requestType: RequestType.delete,
        uri: Uri.parse(ApiConstants.workspacesOperations),
        body: jsonEncode({'workspaceName': workspaceName}));
    if (res.statusCode == 200) {
      return;
    } else {
      throw ServerError();
    }
  }

  Future<List<ServiceAccount>> getServiceAccounts(Workspace workspace) async {
    final uri = Uri.parse(ApiConstants.workspaceAccountOperations)
        .replace(queryParameters: {'workspace': workspace.name});
    final res = await super.authRequest(
        requestType: RequestType.get,
        uri: uri,
        body: jsonEncode({'workspaceName': workspace.name}));
    if (res.statusCode == 200) {
      final parsed = jsonDecode(res.body);
      List<ServiceAccount> accounts = [];
      final accountsJson = parsed['accounts'] as List;
      for (var element in accountsJson) {
        accounts.add(ServiceAccount(
          id: element['id'],
          name: element['name'],
          username: element['username'],
          workspace: workspace
        ));
      }
      return accounts;
    } else {
      throw ServerError();
    }
  }

  Future<void> deleteServiceAccount(
      String workspaceName, String username) async {
    final res = await super.authRequest(
        requestType: RequestType.delete,
        uri: Uri.parse(ApiConstants.workspaceAccountOperations),
        body: jsonEncode({
          'workspaceName': workspaceName,
          'username': username,
        }));
    if (res.statusCode != 200) {
      throw ServerError();
    }
  }

  Future<List<Role>> getRoles(String workspaceName) async {
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

  Future<Role> createRole(String workspaceName, Role role) async {
    final res = await super.authRequest(
      requestType: RequestType.put,
      uri: Uri.parse(ApiConstants.roleOperations),
      body: jsonEncode({
        'workspaceName': workspaceName,
        'name': role.name,
        'description': role.description,
        'canRead': role.canRead,
        'canCreate': role.canCreate,
        'canDelete': role.canDelete,
      }),
    );
    if (res.statusCode == 201) {
      final parsed = jsonDecode(res.body);
      return Role.fromJson(parsed['role']);
    } else {
      throw ServerError();
    }
  }

  Future<Role> updateRole(String workspaceName, Role role) async {
    final res = await super.authRequest(
      requestType: RequestType.patch,
      uri: Uri.parse(ApiConstants.roleOperations),
      body: jsonEncode({
        'workspaceName': workspaceName,
        'id': role.id,
        'name': role.name,
        'description': role.description,
        'canRead': role.canRead,
        'canCreate': role.canCreate,
        'canDelete': role.canDelete,
      }),
    );
    if (res.statusCode == 200) {
      final parsed = jsonDecode(res.body);
      return Role.fromJson(parsed['role']);
    } else {
      throw ServerError();
    }
  }

  Future<void> deleteRole(String workspaceName, Role role) async {
    final res = await super.authRequest(
      requestType: RequestType.delete,
      uri: Uri.parse(ApiConstants.roleOperations),
      body: jsonEncode({
        'workspaceName': workspaceName,
        'roleId': role.id,
      }),
    );
    if (res.statusCode != 200) {
      throw ServerError();
    }
  }

  Future<List<Role>> getAccountRoles(String workspaceName, String accId) async {
    final uri = Uri.parse(ApiConstants.accountRoles).replace(
      queryParameters: { 'workspaceName': workspaceName, 'accountId': accId }
    );
    final res = await authRequest(
      requestType: RequestType.get, uri: uri
    );
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

  Future<void> assignRole(String workspaceName, String accountId, String roleId) async {
    final res = await authRequest(
      requestType: RequestType.post,
      uri: Uri.parse(ApiConstants.accountRoleAssign),
      body: jsonEncode({ 'workspaceName': workspaceName, 'accountId': accountId, 'roleId': roleId })
    );
    if (res.statusCode == 200) {
      return;
    } else {
      throw ServerError();
    }
  }

  Future<void> unAssignRole(String workspaceName, String accountId, String roleId) async {
    final res = await authRequest(
      requestType: RequestType.delete,
      uri: Uri.parse(ApiConstants.accountRoleAssign),
      body: jsonEncode({
        'workspaceName': workspaceName, 'roleId': roleId, 'accountId': accountId,
      })
    );
    if (res.statusCode == 200) {
      return;
    } else {
      throw ServerError();
    }
  }
}
