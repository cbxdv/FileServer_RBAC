import 'package:file_picker/file_picker.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:fs_frontend/blocs/fs/fs_bloc.dart';
import 'package:fs_frontend/blocs/transfer/transfer_bloc.dart';
import 'package:fs_frontend/models/file.dart';
import 'package:fs_frontend/utilities/file_size_utilities.dart';
import 'package:intl/intl.dart';

class FileRow extends StatefulWidget {
  const FileRow({super.key, required this.file});

  final FileModel file;

  @override
  State<FileRow> createState() => _FileRowState();
}

class _FileRowState extends State<FileRow> {
  bool isHovered = false;

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: () {
        showFileDialog(context, widget.file);
      },
      onHover: (value) {
        setState(() {
          isHovered = value;
        });
      },
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 100),
        color: isHovered ? const Color.fromARGB(255, 249, 245, 246) : Colors.white,
        height: 60,
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 30),
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              Container(
                height: 30,
                width: 30,
                decoration: BoxDecoration(
                  color: const Color.fromARGB(255, 188, 232, 242),
                  borderRadius: BorderRadius.circular(5),
                ),
                child: const Icon(
                  Icons.insert_drive_file_outlined,
                  size: 20,
                  color: Color.fromARGB(255, 20, 30, 70),
                ),
              ),
              const SizedBox(width: 30),
              SizedBox(
                width: 300,
                child: Text(
                  widget.file.name,
                  overflow: TextOverflow.ellipsis,
                ),
              ),
              const SizedBox(width: 10),
              SizedBox(
                width: 100,
                child: Text(getSizeString(widget.file.size)),
              ),
              SizedBox(
                width: 100,
                child: Text(DateFormat.yMMMd().format(widget.file.createdOn)),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

void showFileDialog(BuildContext context, FileModel file) {
  showModalBottomSheet(context: context, builder: (_) {
    return SizedBox(
      width: 500,
      child: Padding(
        padding: const EdgeInsets.all(30.0),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(file.name, style: const TextStyle(fontSize: 20, fontWeight: FontWeight.bold, color: Color.fromARGB(255, 20, 30, 70),)),
            const SizedBox(height: 20),
            const Text('Size', style: TextStyle(fontSize: 12, fontWeight: FontWeight.w600)),
            Text(getSizeString(file.size), style: const TextStyle(fontSize: 14)),
            const SizedBox(height: 10),
            const Text('Created On', style: TextStyle(fontSize: 12, fontWeight: FontWeight.w600)),
            Text('${DateFormat.yMd().format(file.createdOn)} - ${DateFormat.jm().format(file.createdOn)}', style: const TextStyle(fontSize: 14)),
            const SizedBox(height: 10),
            const Text('Location', style: TextStyle(fontSize: 12, fontWeight: FontWeight.w600)),
            Text(file.location, style: const TextStyle(fontSize: 14), overflow: TextOverflow.ellipsis),
            const SizedBox(height: 30),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                OutlinedButton(onPressed: () async {
                  String? result = await FilePicker.platform.saveFile(
                    fileName: file.name
                  );
                  if (result == null) {
                    return;
                  }
                  if (context.mounted) {
                    context.read<TransferBloc>().add(
                        DownloadRequestEvent(
                          fileInCloud: file,
                          downloadLocationString: result,
                        )
                    );
                  }
                }, child: const Text('Download')),
                const SizedBox(width: 20),
                OutlinedButton(onPressed: () {
                  Navigator.pop(context);
                  context.read<FSBloc>().add(DeleteFileRequest(location: file.location));
                }, child: const Text('Delete')),
                const SizedBox(width: 20),
                OutlinedButton(onPressed: () {
                  Navigator.pop(context);
                }, child: const Text('Dismiss')),
              ],
            )
          ],
        ),
      ),
    );
  },);
}
