import 'dart:io';

import 'package:file_picker/file_picker.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:fs_frontend/blocs/auth/auth_bloc.dart';
import 'package:fs_frontend/blocs/fs/fs_bloc.dart';
import 'package:fs_frontend/blocs/transfer/transfer_bloc.dart';
import 'package:fs_frontend/models/account.dart';
import 'package:fs_frontend/models/directory.dart';
import 'package:fs_frontend/models/file.dart';
import 'package:fs_frontend/models/role.dart';
import 'package:fs_frontend/models/workspace.dart';
import 'package:fs_frontend/pages/shared_explorer.dart';
import 'package:fs_frontend/widgets/add_role_dialog.dart';
import 'package:fs_frontend/widgets/directory_row.dart';
import 'package:fs_frontend/widgets/explorer_sidebar.dart';
import 'package:fs_frontend/widgets/file_row.dart';
import 'package:fs_frontend/widgets/location_header.dart';

class Explorer extends StatefulWidget {
  const Explorer({super.key, required this.currentWorkspace});

  final Workspace currentWorkspace;

  @override
  State<Explorer> createState() => _ExplorerState();
}

class _ExplorerState extends State<Explorer> {
  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: MultiBlocProvider(
        providers: [
          BlocProvider(
            create: (_) => FSBloc(
              secureStorage: context.read<AuthBloc>().secureStorage,
              workspace: widget.currentWorkspace,
            ),
          ),
          BlocProvider.value(value: context.read<TransferBloc>()),
        ],
        child: ExplorerMain(currentWorkspace: widget.currentWorkspace.name),
      ),
    );
  }
}

class ExplorerMain extends StatefulWidget {
  const ExplorerMain({super.key, required this.currentWorkspace});

  final String currentWorkspace;

  @override
  State<ExplorerMain> createState() => _ExplorerMainState();
}

class _ExplorerMainState extends State<ExplorerMain> {
  @override
  void initState() {
    context
        .read<FSBloc>()
        .add(FSFetchRequest(location: widget.currentWorkspace));
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return BlocListener<FSBloc, FSState>(
      listener: (BuildContext context, FSState state) {
        if (state is FSOperationError) {
          // Navigator.pop(context);
          showDialog(context: context, builder: (context) => AlertDialog(
            title: Text(state.title),
            content: Text(state.description),
            actions: [TextButton(onPressed: () {Navigator.of(context).pop();}, child: const Text('OK'))],
          ),);
        }
      },
      child: const ExplorerBody(),
    );
  }
}

class ExplorerBody extends StatelessWidget {
  const ExplorerBody({super.key});

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        const ExplorerSidebar(),
        Expanded(
          child: Padding(
            padding: const EdgeInsets.all(60.0),
            child: BlocBuilder<FSBloc, FSState>(
              buildWhen: (previous, current) => current is! FSActionState,
              builder: (context, state) {
                if (state is FSFetching || state is FSOperationLoadingState) {
                  return const Center(child: CircularProgressIndicator());
                }
                if (state is FSPermissionDeniedState) {
                  final workspaceName = state.location.split("/").first;
                  return Center(
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        const Text(
                          'Permission Denied',
                          style: TextStyle(
                            fontSize: 24,
                            fontWeight: FontWeight.bold,
                            color: Color.fromARGB(255, 20, 30, 70),
                          ),
                        ),
                        const SizedBox(height: 20),
                        OutlinedButton(onPressed: () {
                          Navigator.of(context).push(MaterialPageRoute(
                            builder: (_) => MultiBlocProvider(
                                providers: [
                                  BlocProvider.value(value: context.read<FSBloc>()),
                                  BlocProvider.value(value: context.read<TransferBloc>()),
                                ],
                                child: SharedExplorer(workspaceName: workspaceName)
                            )));
                        }, child: const Text('Shared'))
                      ],
                    ),
                  );
                }
                if (state is FSFetched) {
                  final data = state.dir;
                  return Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Row(
                        // mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        // crossAxisAlignment: CrossAxisAlignment.center,
                        children: [
                          LocationHeader(location: data.location),
                          Container(
                            margin: const EdgeInsets.only(left: 10),
                            color: Colors.white,
                            child: Row(
                              children: [
                                IconButton(onPressed: () {
                                  context.read<FSBloc>().add(FSFetchRequest(
                                      location: data.location
                                  ));
                                }, icon: const Icon(Icons.refresh)),
                                const SizedBox(width: 20),
                                const AddNewButton(),
                              ],
                            ),
                          )
                        ],
                      ),
                      const SizedBox(height: 30),
                      const HeaderRow(),
                      Expanded(
                        child: Builder(
                          builder: (_) {
                            if (data.contents.isEmpty) {
                              return const Center(child: Text('No data found'),);
                            }
                            return ListView.builder(
                              itemCount: data.contents.length,
                              itemBuilder: ((context, index) {
                                // return Row();
                                if (data.contents[index] is DirectoryModel) {
                                  return DirectoryRow(
                                    directory: data.contents[index],
                                  );
                                } else if (data.contents[index] is FileModel) {
                                  return FileRow(file: data.contents[index]);
                                }
                                return Container();
                              }),
                            );
                          }
                        ),
                      ),
                    ],
                  );
                }
                return Container();
              },
            ),
          ),
        ),
      ],
    );
  }
}

enum AddNewEnum { newFolder, fileUpload, role }

class HeaderRow extends StatelessWidget {
  const HeaderRow({super.key});

  @override
  Widget build(BuildContext context) {
    return Container(
      color: Colors.white,
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
                // color: Color.fromARGB(255, 255, 220, 170),
                borderRadius: BorderRadius.circular(5),
              ),
              child: Container(),
            ),
            const SizedBox(width: 30),
            const SizedBox(
              width: 300,
              child: Text(
                'Name',
                overflow: TextOverflow.ellipsis,
                style: TextStyle(color: Color.fromARGB(255, 159, 159, 159)),
              ),
            ),
            const SizedBox(width: 10),
            const SizedBox(
              width: 100,
              child: Text(
                'Size',
                style: TextStyle(color: Color.fromARGB(255, 159, 159, 159)),
              ),
            ),
            const SizedBox(
              width: 100,
              child: Text(
                'Created on',
                style: TextStyle(color: Color.fromARGB(255, 159, 159, 159)),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class AddNewButton extends StatelessWidget {
  const AddNewButton({super.key});

  uploadFile(BuildContext context) async {
    FilePickerResult? result = await FilePicker.platform.pickFiles(
      dialogTitle: 'Select file to upload',
      allowMultiple: false,
    );
    if (result == null) {
      return;
    }
    File selectedFile = File(result.files.single.path!);
    if (context.mounted && context.read<FSBloc>().state is FSFetched) {
      final state = context.read<FSBloc>().state as FSFetched;

      if (context.mounted) {
        context.read<TransferBloc>().add(
          UploadRequestEvent(file: selectedFile, location: state.dir.location)
      );
      }
    }
  }

  addRoleHandler(BuildContext context) async {
    if (context.read<FSBloc>().state is! FSFetched) {
      return;
    }
    final dir = (context.read<FSBloc>().state as FSFetched).dir;
    final dirDetails = await context.read<FSBloc>().fsRepo.getDirectoryDetails(dir.location);
    final workspaceName = dirDetails.location.split("/").first;
    Role? selectedRole;
    if (context.mounted) {
      selectedRole = await showAddRoleDialog(
          context, dirDetails.roles,
              () => context.read<FSBloc>().fsRepo.getAllRoles(workspaceName)
      );
    }
    if (selectedRole == null) {
      return;
    }
    try {
      if (context.mounted) {
        context.read<FSBloc>().fsRepo.addRoleToFS(
            dir.location, selectedRole
        );
      }
    } catch (_) {
      // Do nothing
    }
  }

  @override
  Widget build(BuildContext context) {
    final dirLocation =
        (context.read<FSBloc>().state as FSFetched).dir.location;

    createDir(String name) {
      Navigator.pop(context);
      if (context.read<FSBloc>().state is FSFetching) {
        return;
      }
      context
          .read<FSBloc>()
          .add(CreateDirRequest(newDirName: name, parentLocation: dirLocation));
    }

    return PopupMenuButton<AddNewEnum>(
      color: const Color.fromARGB(255, 249, 245, 246),
      surfaceTintColor: const Color.fromARGB(255, 249, 245, 246),
      position: PopupMenuPosition.under,
      offset: const Offset(0, 5),
      itemBuilder: (context) => [
        const PopupMenuItem(
          value: AddNewEnum.newFolder,
          child: Row(children: [
            Icon(Icons.folder_outlined),
            SizedBox(width: 10),
            Text('Folder')
          ]),
        ),
        const PopupMenuItem(
            value: AddNewEnum.fileUpload,
            child: Row(children: [
              Icon(Icons.upload_file_outlined),
              SizedBox(width: 10),
              Text('Upload file')
            ])),
        if ((context.read<AuthBloc>().state as AuthLoggedInState).account is OwnerAccount) ...[
          const PopupMenuItem(
              value: AddNewEnum.role,
              child: Row(children: [
                Icon(Icons.grid_3x3),
                SizedBox(width: 10),
                Text('Role')
              ]))
        ],
      ],
      onSelected: (value) => {
        switch (value) {
          AddNewEnum.newFolder => showCreateFolderButton(context, createDir),
          AddNewEnum.fileUpload => uploadFile(context),
          AddNewEnum.role => addRoleHandler(context),
        }
      },
      child: Container(
        height: 36,
        width: 150,
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(10),
          color: const Color.fromARGB(255, 20, 30, 70),
        ),
        child: const Center(
          child: Text(
            'Add new',
            style: TextStyle(
              color: Color.fromARGB(255, 244, 238, 238),
              fontWeight: FontWeight.bold,
            ),
          ),
        ),
      ),
    );
  }
}

void showCreateFolderButton(BuildContext context, Function createDir) {
  final TextEditingController controller = TextEditingController();

  showDialog(
    context: context,
    builder: (context) => AlertDialog(
      title: const Text(
        'Enter a name for the folder',
        style: TextStyle(fontSize: 20),
      ),
      content: SizedBox(
        width: 200,
        child: TextField(
          controller: controller,
          decoration: const InputDecoration(hintText: 'New directory name'),
        ),
      ),
      actions: [
        TextButton(
          onPressed: () {
            Navigator.pop(context);
          },
          child: const Text('Discard'),
        ),
        TextButton(
            onPressed: () {
              createDir(controller.text);
            },
            child: const Text('Create folder')),
      ],
    ),
  );
}
