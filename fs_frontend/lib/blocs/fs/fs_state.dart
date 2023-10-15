part of 'fs_bloc.dart';

abstract class FSState {}

abstract class FSActionState extends FSState {}

class FSInitial extends FSState {}

class FSFetching extends FSState {}

class FSFetched extends FSState {
  final DirectoryModel dir;
  FSFetched({required this.dir});
}

class FSError extends FSState {}

class FSOperationLoadingState extends FSState {}

class FSPermissionDeniedState extends FSState {
  final String location;

  FSPermissionDeniedState({required this.location});
}

class FSOperationError extends FSActionState {
  final String title;
  final String description;

  FSOperationError({required this.title, required this.description});
}