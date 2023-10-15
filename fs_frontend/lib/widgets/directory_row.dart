
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:fs_frontend/blocs/auth/auth_bloc.dart';
import 'package:fs_frontend/blocs/fs/fs_bloc.dart';
import 'package:fs_frontend/models/account.dart';
import 'package:fs_frontend/models/directory.dart';
import 'package:fs_frontend/widgets/directory_details.dart';
import 'package:intl/intl.dart';

class DirectoryRow extends StatefulWidget {
  const DirectoryRow({super.key, required this.directory});

  final DirectoryModel directory;

  @override
  State<DirectoryRow> createState() => _DirectoryRowState();
}

class _DirectoryRowState extends State<DirectoryRow> {
  bool isHovered = false;

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: () {showDirectoryDialog(context, widget.directory);},
      onDoubleTap: () {
        context.read<FSBloc>().add(FSLocationChangeRequest(location: widget.directory.location));
        if (
        (context.read<FSBloc>().state is FSPermissionDeniedState) &&
        (context.read<AuthBloc>().state as AuthLoggedInState).account is ServiceAccount) {
          Navigator.of(context).pop();
        }
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
                  color: const Color.fromARGB(255, 255, 220, 170),
                  borderRadius: BorderRadius.circular(5),
                ),
                child: const Icon(
                  Icons.folder_outlined,
                  size: 20,
                  color: Color.fromARGB(255, 20, 30, 70),
                ),
              ),
              const SizedBox(width: 30),
              SizedBox(
                width: 300,
                child: Text(
                  widget.directory.name,
                  overflow: TextOverflow.ellipsis,
                ),
              ),
              const SizedBox(width: 10),
              const SizedBox(
                width: 100,
              ),
              SizedBox(
                width: 100,
                child:
                Text(DateFormat.yMMMd().format(widget.directory.createdOn)),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

void showDirectoryDialog(BuildContext context, DirectoryModel dir) {
  showModalBottomSheet(context: context, builder: (_) {
    return SizedBox(
      width: 500,
      child: BlocProvider.value(
        value: context.read<FSBloc>(),
        child: DirectoryDetails(dir: dir),
      ),
    );
  },);
}

