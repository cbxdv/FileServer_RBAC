part of 'transfer_bloc.dart';

abstract class TransferEvent {}

class UploadRequestEvent extends TransferEvent {
  final File file;
  final String location;
  UploadRequestEvent({required this.file, required this.location});
}

class DownloadRequestEvent extends TransferEvent {
  final FileModel fileInCloud;
  final String downloadLocationString;
  DownloadRequestEvent({required this.fileInCloud, required this.downloadLocationString});
}

class ClearCompleted extends TransferEvent {}
