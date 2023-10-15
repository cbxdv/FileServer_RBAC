part of 'workspace_bloc.dart';

abstract class WorkspaceEvent {}

class WorkspaceFetchRequest extends WorkspaceEvent {}

class WorkspaceCreate extends WorkspaceEvent {
  final String newWorkspaceName;
  WorkspaceCreate({required this.newWorkspaceName});
}

class WorkspaceDelete extends WorkspaceEvent {
  final String workspaceName;
  WorkspaceDelete({required this.workspaceName});
}
