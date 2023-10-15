import 'dart:async';
import 'dart:io';

import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:fs_frontend/models/file.dart';
import 'package:fs_frontend/models/transfer_task.dart';
import 'package:fs_frontend/repos/transfer_repo.dart';

part 'transfer_event.dart';
part 'transfer_state.dart';

class TransferBloc extends Bloc<TransferEvent, TransferState> {
  late TransferRepo transferRepo;

  TransferBloc({ required this.transferRepo }) : super(TransferState(transfers: [])) {
    on<UploadRequestEvent>(upload);
    on<DownloadRequestEvent>(download);
    on<ClearCompleted>(clearCompleted);
  }

  FutureOr<void> upload(UploadRequestEvent event, Emitter<TransferState> emit) async {
    final transfers = state.transfers;
    final size = await event.file.length();
    final uploadTask = UploadTask(file: event.file, parentLocation: event.location, size: size);
    transfers.add(uploadTask);
    transferRepo.sink.add(uploadTask);
    emit(TransferState(transfers: transfers));
  }

  FutureOr<void> download(DownloadRequestEvent event, Emitter<TransferState> emit) async {
    final transfers = state.transfers;
    List<String> locationSplit = event.downloadLocationString.split('/');
    if (locationSplit.length == 1) {
      locationSplit = event.downloadLocationString.split('\\');
    }
    final downloadTask = DownloadTask(
      fileName: locationSplit.last,
      size: event.fileInCloud.size,
      downloadPathString: event.downloadLocationString,
      fileInCloud: event.fileInCloud
    );
    transfers.add(downloadTask);
    transferRepo.sink.add(downloadTask);
    emit(TransferState(transfers: transfers));
  }

  FutureOr<void> clearCompleted(ClearCompleted event, Emitter<TransferState> emit) async {
    final transfers = state.transfers;
    transfers.removeWhere((element) => element.hasCompleted || element.hasFailed);
    emit(TransferState(transfers: transfers));
  }
}