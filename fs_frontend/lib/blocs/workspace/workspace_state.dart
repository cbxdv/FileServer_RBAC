part of 'workspace_bloc.dart';

abstract class WorkspaceState {}

abstract class WorkspaceActionState extends WorkspaceState {}

class WorkspaceInitial extends WorkspaceState {}

class WorkspaceLoading extends WorkspaceState {}

class WorkspacesLoaded extends WorkspaceState {
  final List<Workspace> workspaces;
  WorkspacesLoaded({required this.workspaces});
}

class WorkspacesFailed extends WorkspaceState {}

class WorkspaceAlreadyExistsState extends WorkspaceActionState {}