import 'dart:async';

import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:fs_frontend/exceptions/exceptions.dart';
import 'package:fs_frontend/models/workspace.dart';
import 'package:fs_frontend/repos/workspace_repo.dart';

part 'workspace_event.dart';
part 'workspace_state.dart';

class WorkspaceBloc extends Bloc<WorkspaceEvent, WorkspaceState> {
  late final WorkspaceRepo wsRepo;
  final FlutterSecureStorage secureStorage;

  WorkspaceBloc({ required this.secureStorage }) : super(WorkspaceInitial()) {
    wsRepo = WorkspaceRepo(secureStorage: secureStorage);
    on<WorkspaceFetchRequest>(fetchWorkspace);
    on<WorkspaceCreate>(createWorkspace);
    on<WorkspaceDelete>(deleteWorkspace);
  }

  FutureOr<void> fetchWorkspace(WorkspaceFetchRequest event, Emitter<WorkspaceState> emit) async {
    emit(WorkspaceLoading());
    try {
      List<Workspace> workspaces = await wsRepo.getAllWorkspaces();
      emit(WorkspacesLoaded(workspaces: workspaces));
    } catch (_) {
      emit(WorkspacesFailed());
    }
  }

  FutureOr<void> createWorkspace(WorkspaceCreate event, Emitter<WorkspaceState> emit) async {
    List<Workspace> workspaces = [];
    if (state is WorkspacesLoaded) {
      workspaces = (state as WorkspacesLoaded).workspaces;
    }
    emit(WorkspaceLoading());
    try {
      Workspace newWorkspace = await wsRepo.createNewWorkspace(event.newWorkspaceName);
      workspaces.add(newWorkspace);
      emit(WorkspacesLoaded(workspaces: workspaces));
    } on WorkspaceAlreadyExists {
      emit(WorkspaceAlreadyExistsState());
      emit(WorkspacesLoaded(workspaces: workspaces));
    } catch (_) {
      emit(WorkspacesFailed());
    }
  }

  FutureOr<void> deleteWorkspace(WorkspaceDelete event, Emitter<WorkspaceState> emit) async {
    List<Workspace> workspaces = [];
    if (state is WorkspacesLoaded) {
      workspaces = (state as WorkspacesLoaded).workspaces;
    }
    emit(WorkspaceLoading());
    try {
      await wsRepo.deleteWorkspace(event.workspaceName);
      workspaces.removeWhere((element) => element.name == event.workspaceName);
      emit(WorkspacesLoaded(workspaces: workspaces));
    } catch (_) {
      emit(WorkspacesLoaded(workspaces: workspaces));
    }
  }
}