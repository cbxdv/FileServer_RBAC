import 'dart:async';

import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:fs_frontend/exceptions/exceptions.dart';
import 'package:fs_frontend/models/directory.dart';
import 'package:fs_frontend/models/file.dart';
import 'package:fs_frontend/models/workspace.dart';
import 'package:fs_frontend/repos/fs_repo.dart';

part 'fs_event.dart';
part 'fs_state.dart';

class FSBloc extends Bloc<FSEvent, FSState> {
  final Workspace workspace;
  late FSRepo fsRepo;
  final FlutterSecureStorage secureStorage;

  FSBloc({required this.secureStorage, required this.workspace})
      : super(FSInitial()) {
    fsRepo = FSRepo(secureStorage: secureStorage);

    on<FSFetchRequest>(fetch);
    on<CreateDirRequest>(createDir);
    on<DeleteDirRequest>(deleteDir);
    on<FSLocationChangeRequest>(changeLocation);
    on<DeleteFileRequest>(deleteFile);
    on<CreateFile>(createFile);
  }

  FutureOr<void> fetch(FSFetchRequest event, Emitter<FSState> emit) async {
    emit(FSFetching());
    try {
      DirectoryModel dir = await fsRepo.getDirectory(event.location);
      emit(FSFetched(dir: dir));
    } on PermissionDenied {
      emit(FSPermissionDeniedState(location: event.location));
    } catch (_) {
      emit(FSError());
    }
  }

  FutureOr<void> createDir(
      CreateDirRequest event, Emitter<FSState> emit) async {
    final DirectoryModel dir = (state as FSFetched).dir;
    emit(FSOperationLoadingState());
    try {
      final newDir =
          await fsRepo.createDirectory(event.newDirName, event.parentLocation);
      dir.contents.add(newDir);
    } on ResourceAlreadyExists catch (e) {
      emit(FSOperationError(
          title: 'Directory already exists', description: e.description));
    } catch (_) {
      emit(FSOperationError(
          title: 'Something unexpected happened',
          description: 'Try again later'));
    } finally {
      emit(FSFetched(dir: dir));
    }
  }

  FutureOr<void> deleteDir(
      DeleteDirRequest event, Emitter<FSState> emit) async {
    final DirectoryModel dir = (state as FSFetched).dir;
    emit(FSOperationLoadingState());
    try {
      await fsRepo.deleteDirectory(event.location);
      dir.contents.removeWhere((element) => element.location == event.location);
    } on UnauthorizedError {
      emit(FSOperationError(
          title: 'Unauthorized', description: 'You are not allowed to delete'));
    } on ChildrenExists catch (e) {
      emit(FSOperationError(
          title: 'Error deleting', description: e.description));
    } catch (_) {
      emit(FSOperationError(
          title: 'Something unexpected happened',
          description: 'Try again later'));
    } finally {
      emit(FSFetched(dir: dir));
    }
  }

  FutureOr<void> changeLocation(
      FSLocationChangeRequest event, Emitter<FSState> emit) async {
    emit(FSFetching());
    try {
      DirectoryModel dir = await fsRepo.getDirectory(event.location);
      emit(FSFetched(dir: dir));
    } on PermissionDenied {
      emit(FSPermissionDeniedState(location: event.location));
    } catch (_) {
      emit(FSError());
    }
  }

  FutureOr<void> deleteFile(
      DeleteFileRequest event, Emitter<FSState> emit) async {
    final DirectoryModel dir = (state as FSFetched).dir;
    emit(FSOperationLoadingState());
    try {
      await fsRepo.deleteFile(event.location);
      dir.contents.removeWhere((element) => element.location == event.location);
    } on UnauthorizedError {
      emit(FSOperationError(
          title: 'Unauthorized', description: 'You are not allowed to delete'));
    } catch (_) {
      emit(FSOperationError(
          title: 'Something unexpected happened',
          description: 'Try again later'));
    } finally {
      emit(FSFetched(dir: dir));
    }
  }

  FutureOr<void> createFile(CreateFile event, Emitter<FSState> emit) {
    final DirectoryModel dir = (state as FSFetched).dir;
    emit(FSOperationLoadingState());
    final flSplit = event.file.location.split('/');
    final parentLocation = flSplit.sublist(0, flSplit.length - 1).join('/');
    if (dir.location == parentLocation) {
      dir.contents.add(event.file);
    }
    emit(FSFetched(dir: dir));
  }
}
