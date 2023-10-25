import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:fs_frontend/blocs/fs/fs_bloc.dart';
import 'package:fs_frontend/models/directory.dart';
import 'package:fs_frontend/models/role.dart';
import 'package:fs_frontend/widgets/add_role_dialog.dart';
import 'package:intl/intl.dart';

class DirectoryDetails extends StatefulWidget {
  const DirectoryDetails({super.key, required this.dir});

  final DirectoryModel dir;

  @override
  State<DirectoryDetails> createState() => _DirectoryDetailsState();
}

class _DirectoryDetailsState extends State<DirectoryDetails> {
  bool isLoading = true;
  String err = "";
  late DirectoryModel dirDetails;

  void fetchData() async {
    setState(() {
      isLoading = true;
    });
    try {
      final dirRes = await context
          .read<FSBloc>()
          .fsRepo
          .getDirectoryDetails(widget.dir.location);
      setState(() {
        dirDetails = dirRes;
        isLoading = false;
      });
    } catch (_) {
      setState(() {
        isLoading = false;
        err = "Server error";
      });
    }
  }

  void removeRoleHandler(Role role) async {
    try {
      context.read<FSBloc>().fsRepo.removeRoleFS(dirDetails.location, role);
      setState(() {
        dirDetails.roles.remove(role);
      });
    } catch (_) {
      setState(() {
        isLoading = false;
        err = "Server error";
      });
    }
  }

  void addRoleHandler() async {
    final workspaceName = dirDetails.location.split("/").first;
    final roleSelected = await showAddRoleDialog(context, dirDetails.roles,
        () => context.read<FSBloc>().fsRepo.getAllRoles(workspaceName));
    if (roleSelected == null) {
      return;
    }
    try {
      setState(() {
        isLoading = true;
      });
      if (context.mounted) {
        await context
            .read<FSBloc>()
            .fsRepo
            .addRoleToFS(widget.dir.location, roleSelected);
      }
      setState(() {
        isLoading = false;
        dirDetails.roles.add(roleSelected);
      });
    } catch (err) {
      setState(() {
        isLoading = false;
      });
    }
  }

  @override
  void initState() {
    dirDetails = widget.dir;
    fetchData();
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    if (isLoading) {
      return const SizedBox(
          width: 30,
          height: 120,
          child: Center(child: CircularProgressIndicator()));
    }
    return Padding(
      padding: const EdgeInsets.all(30.0),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(dirDetails.name,
              style: const TextStyle(
                fontSize: 20,
                fontWeight: FontWeight.bold,
                color: Color.fromARGB(255, 20, 30, 70),
              )),
          const SizedBox(height: 20),
          const Text('Created On',
              style: TextStyle(fontSize: 12, fontWeight: FontWeight.w600)),
          Text(
              '${DateFormat.yMd().format(dirDetails.createdOn)} - ${DateFormat.jm().format(dirDetails.createdOn)}',
              style: const TextStyle(fontSize: 14)),
          const SizedBox(height: 10),
          const Text('Location',
              style: TextStyle(fontSize: 12, fontWeight: FontWeight.w600)),
          Text(dirDetails.location,
              style: const TextStyle(fontSize: 14),
              overflow: TextOverflow.ellipsis),
          const SizedBox(height: 20),
          if (dirDetails.roles.isNotEmpty) ...[
            const Text('Roles',
                style: TextStyle(fontSize: 12, fontWeight: FontWeight.w600)),
            const SizedBox(height: 5),
            Wrap(
              children: List.generate(
                dirDetails.roles.length,
                (index) => Container(
                  margin: const EdgeInsets.symmetric(horizontal: 5),
                  child: Chip(
                    backgroundColor: Colors.transparent,
                    side: BorderSide(color: Colors.grey.shade200),
                    elevation: 0,
                    padding: EdgeInsets.zero,
                    materialTapTargetSize: MaterialTapTargetSize.shrinkWrap,
                    label: Text(dirDetails.roles[index].name),
                    deleteIcon: const Icon(
                      Icons.remove_circle_outline,
                      size: 16,
                    ),
                    deleteButtonTooltipMessage: 'Remove',
                    onDeleted: () => removeRoleHandler(dirDetails.roles[index]),
                  ),
                ),
              ),
            )
          ],
          const SizedBox(height: 30),
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              OutlinedButton(
                  onPressed: addRoleHandler, child: const Text('Add Role')),
              const SizedBox(width: 20),
              OutlinedButton(
                  onPressed: () {
                    Navigator.pop(context);
                    context
                        .read<FSBloc>()
                        .add(DeleteDirRequest(location: dirDetails.location));
                  },
                  child: const Text('Delete')),
              const SizedBox(width: 20),
              OutlinedButton(
                  onPressed: () {
                    Navigator.pop(context);
                  },
                  child: const Text('Dismiss')),
            ],
          )
        ],
      ),
    );
  }
}
