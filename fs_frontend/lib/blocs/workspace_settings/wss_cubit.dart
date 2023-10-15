import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:fs_frontend/blocs/workspace_settings/wss_state.dart';
import 'package:fs_frontend/models/role.dart';
import 'package:fs_frontend/models/workspace.dart';
import 'package:fs_frontend/repos/workspace_repo.dart';

class WSSLoading extends WSSState {}

class WSSServerErr extends WSSState {}

class WSSCubit extends Cubit<WSSState> {
  final WorkspaceRepo workspaceRepo;
  final Workspace workspace;

  WSSCubit({required this.workspaceRepo, required this.workspace})
      : super(WSSState());

  Future<void> updateWorkspaceName(String name) async {
    final oldState = state;
    emit(WSSLoading());
    try {
      await workspaceRepo.updateWorkspaceName(workspace.name, name);
    } catch (_) {
      emit(WSSServerErr());
      emit(oldState);
    }
  }

  Future<void> getServiceAccounts() async {
    final newState = WSSState.copyFrom(state);
    newState.isLoadingAccounts = true;
    emit(newState);
    try {
      final accounts = await workspaceRepo.getServiceAccounts(workspace);
      final newState = WSSState.copyFrom(state);
      newState.accounts = accounts;
      newState.isLoadingAccounts = false;
      emit(newState);
    } catch (_) {
      emit(WSSServerErr());
      emit(WSSState());
    }
  }

  Future<void> createServiceAccount(
      String name, String username, String password) async {
    final oldState = WSSState.copyFrom(state);
    oldState.isLoadingAccounts = true;
    emit(oldState);
    try {
      final newAcc = await workspaceRepo.createServiceAccount(
          workspace, name, username, password);
      final newState = WSSState.copyFrom(oldState);
      newState.isLoadingAccounts = false;
      newState.accounts.add(newAcc);
      emit(newState);
    } catch (_) {
      emit(WSSServerErr());
      oldState.isLoadingAccounts = false;
      emit(oldState);
    }
  }

  Future<void> deleteServiceAccount(String username) async {
    final oldState = state;
    emit(WSSLoading());
    try {
      oldState.accounts.removeWhere((element) => element.username == username);
      await workspaceRepo.deleteServiceAccount(workspace.name, username);
      emit(oldState);
    } catch (_) {
      emit(WSSServerErr());
      emit(oldState);
    }
  }

  Future<void> getRoles() async {
    final oldState = WSSState();
    oldState.isLoadingRoles = true;
    emit(oldState);
    try {
      final roles = await workspaceRepo.getRoles(workspace.name);
      final newState = WSSState.copyFrom(oldState);
      newState.isLoadingRoles = false;
      newState.roles = roles;
      emit(newState);
    } catch (_) {
      final s = WSSState.copyFrom(oldState);
      s.isLoadingRoles = false;
      emit(WSSServerErr());
      emit(s);
    }
  }

  Future<void> createRole({
    required String name,
    required String description,
    required bool canRead,
    required bool canCreate,
    required bool canDelete,
  }) async {
    final oldState = WSSState.copyFrom(state);
    oldState.isLoadingRoles = true;
    emit(oldState);
    final role = Role(
        id: 'id',
        name: name,
        description: description,
        canRead: canRead,
        canCreate: canCreate,
        canDelete: canDelete);
    try {
      final newRole = await workspaceRepo.createRole(workspace.name, role);
      final newState = WSSState.copyFrom(oldState);
      newState.roles.add(newRole);
      newState.isLoadingRoles = false;
      emit(newState);
    } catch (_) {
      emit(WSSServerErr());
      emit(oldState);
    }
  }

  Future<void> updateRole({
    required String id,
    required String name,
    required String description,
    required bool canRead,
    required bool canCreate,
    required bool canDelete,
  }) async {
    final oldState = WSSState.copyFrom(state);
    oldState.isLoadingRoles = true;
    emit(oldState);
    final role = Role(
      id: id,
      name: name,
      description: description,
      canRead: canRead,
      canCreate: canCreate,
      canDelete: canDelete,
    );
    try {
      final newRole = await workspaceRepo.updateRole(workspace.name, role);
      final newState = WSSState.copyFrom(oldState);
      newState.roles.removeWhere((element) => element.id == role.id);
      newState.roles.add(newRole);
      emit(newState);
    } catch (_) {
      emit(WSSServerErr());
      oldState.isLoadingRoles = false;
      emit(oldState);
    }
  }

  Future<void> deleteRole(Role role) async {
    final oldState = WSSState.copyFrom(state);
    oldState.isLoadingRoles = true;
    emit(oldState);
    try {
      await workspaceRepo.deleteRole(workspace.name, role);
      final roles = oldState.roles;
      roles.removeWhere((element) => element.id == role.id);
      final newState = WSSState.copyFrom(oldState);
      newState.roles = roles;
      newState.isLoadingRoles = false;
      emit(newState);
    } catch (_) {
      emit(WSSServerErr());
      emit(oldState);
    }
  }
}
