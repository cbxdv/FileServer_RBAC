import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:fs_frontend/blocs/auth/auth_bloc.dart';
import 'package:fs_frontend/blocs/transfer/transfer_bloc.dart';
import 'package:fs_frontend/blocs/workspace/workspace_bloc.dart';
import 'package:fs_frontend/pages/account_settings.dart';
import 'package:fs_frontend/pages/explorer.dart';
import 'package:fs_frontend/repos/transfer_repo.dart';
import 'package:fs_frontend/widgets/transfer_pane.dart';

class WorkspaceSelect extends StatefulWidget {
  const WorkspaceSelect({super.key});

  @override
  State<WorkspaceSelect> createState() => _WorkspaceSelectState();
}

class _WorkspaceSelectState extends State<WorkspaceSelect> {
  @override
  void initState() {
    context.read<WorkspaceBloc>().add(WorkspaceFetchRequest());
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
        create: (_) => TransferBloc(
                transferRepo: TransferRepo(
              secureStorage: context.read<AuthBloc>().secureStorage,
            )),
        child: const WorkspaceSelectMain());
  }
}

class WorkspaceSelectMain extends StatefulWidget {
  const WorkspaceSelectMain({super.key});

  @override
  State<WorkspaceSelectMain> createState() => _WorkspaceSelectMainState();
}

class _WorkspaceSelectMainState extends State<WorkspaceSelectMain> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          showAddDialog(context, context.read<WorkspaceBloc>());
        },
        child: const Icon(Icons.add),
      ),
      body: Row(
        children: [
          const WorkspaceSelectPageSidebar(),
          Expanded(
            child: Container(
              margin: const EdgeInsets.all(30),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.start,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 15),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        const Text(
                          'Select Workspace',
                          style: TextStyle(
                            fontSize: 24,
                            fontWeight: FontWeight.bold,
                              color: Color.fromARGB(255, 20, 30, 70),
                          ),
                        ),
                        IconButton(
                            onPressed: () {
                              context
                                  .read<WorkspaceBloc>()
                                  .add(WorkspaceFetchRequest());
                            },
                            icon: const Icon(Icons.refresh))
                      ],
                    ),
                  ),
                  const SizedBox(height: 20),
                  Expanded(
                    child: BlocConsumer<WorkspaceBloc, WorkspaceState>(
                      buildWhen: (previous, current) =>
                          current is! WorkspaceActionState,
                      listener: (context, state) {
                        if (state is WorkspaceAlreadyExistsState) {
                          showAlreadyExistsDialog(context);
                        }
                      },
                      builder: (context, state) {
                        if (state is WorkspacesFailed) {
                          return const Center(
                              child: Text('Loading workspaces failed'));
                        }
                        if (state is WorkspaceLoading) {
                          return const Center(child: CircularProgressIndicator());
                        }
                        if (state is WorkspacesLoaded) {
                          return const WorkspaceBody();
                        }
                        return Container();
                      },
                    ),
                  )
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class WorkspaceSelectPageSidebar extends StatelessWidget {
  const WorkspaceSelectPageSidebar({super.key});

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: 80,
      child: Container(
        decoration: const BoxDecoration(
          borderRadius: BorderRadius.only(
            topRight: Radius.circular(10),
            bottomRight: Radius.circular(10),
          ),
          color: Color.fromARGB(255, 244, 238, 238),
        ),
        child: Column(
          children: [
            const SizedBox(height: 60),
            const Icon(
              Icons.cloud_circle_outlined,
              size: 40,
            ),
            Expanded(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  IconButton(
                    onPressed: () {
                      showTransferPane(context);
                    },
                    icon: const Icon(Icons.downloading),
                    tooltip: 'Transfers',
                    iconSize: 30,
                  ),
                  const SizedBox(height: 30),
                  IconButton(
                    onPressed: () {
                      Navigator.push(
                          context,
                          MaterialPageRoute(
                            builder: (context) => const AccountSettings(),
                          ));
                    },
                    icon: const Icon(Icons.settings_outlined),
                    tooltip: 'Account Settings',
                    iconSize: 30,
                  ),
                ],
              ),
            ),
            IconButton(
              onPressed: () {
                context.read<AuthBloc>().add(LogoutEvent());
              },
              tooltip: 'Logout',
              icon: const Icon(Icons.logout),
              iconSize: 30,
            ),
            const SizedBox(height: 30),
          ],
        ),
      ),
    );
  }
}

void showAddDialog(BuildContext context, WorkspaceBloc wsBloc) {
  showDialog(
      context: context,
      builder: (context) {
        final newWsController = TextEditingController();
        return AlertDialog(
          title: const Text("New Workspace"),
          content: TextField(
            controller: newWsController,
            decoration: const InputDecoration(
              hintText: 'Workspace Name',
            ),
          ),
          actions: [
            TextButton(
                onPressed: () {
                  Navigator.pop(context);
                },
                child: const Text('Discard')),
            TextButton(
                onPressed: () {
                  if (newWsController.text.isEmpty) {
                    return;
                  }
                  wsBloc.add(
                      WorkspaceCreate(newWorkspaceName: newWsController.text));
                  Navigator.pop(context);
                },
                child: const Text('Create')),
          ],
        );
      });
}

void showAlreadyExistsDialog(BuildContext context) {
  showDialog(
    context: context,
    builder: (context) => AlertDialog(
      title: const Text("Workspace exists"),
      content: const Text('Try a different workspace name'),
      actions: [
        TextButton(
            onPressed: () {
              Navigator.of(context).pop();
            },
            child: const Text('OK'))
      ],
    ),
  );
}

class WorkspaceBody extends StatefulWidget {
  const WorkspaceBody({super.key});

  @override
  State<WorkspaceBody> createState() => _WorkspaceBodyState();
}

class _WorkspaceBodyState extends State<WorkspaceBody> {
  @override
  Widget build(BuildContext context) {
    final workspaces =
        (context.read<WorkspaceBloc>().state as WorkspacesLoaded).workspaces;
    if (workspaces.isEmpty) {
      return const Center(child: Text('No workspaces found'));
    }
    return ListView.builder(
      itemCount: workspaces.length,
      itemBuilder: (_, index) {
        return Card(
          shadowColor: Colors.transparent,
          child: ListTile(
            title: Text(workspaces[index].name),
            onTap: () {
              Navigator.of(context).push(MaterialPageRoute(
                builder: (c) => MultiBlocProvider(
                  providers: [
                    BlocProvider.value(value: context.read<WorkspaceBloc>()),
                    BlocProvider.value(value: context.read<TransferBloc>()),
                  ],
                  child: Explorer(currentWorkspace: workspaces[index]),
                ),
              ));
            },
          ),
        );
      },
    );
  }
}
