part of 'fs_bloc.dart';

abstract class FSEvent {}

class FSFetchRequest extends FSEvent {
  final String location;
  FSFetchRequest({required this.location});
}

class FSLocationChangeRequest extends FSEvent {
  final String location;
  FSLocationChangeRequest({required this.location});
}

class CreateDirRequest extends FSEvent {
  final String newDirName;
  final String parentLocation;
  CreateDirRequest({required this.newDirName, required this.parentLocation});
}

class DeleteDirRequest extends FSEvent {
  final String location;
  DeleteDirRequest({required this.location});
}

class DeleteFileRequest extends FSEvent {
  final String location;
  DeleteFileRequest({required this.location});
}

class CreateFile extends FSEvent {
  final FileModel file;
  CreateFile({required this.file});
}