import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:fs_frontend/blocs/workspace/workspace_bloc.dart';
import 'package:fs_frontend/blocs/workspace_settings/wss_cubit.dart';
import 'package:fs_frontend/blocs/workspace_settings/wss_state.dart';
import 'package:fs_frontend/models/workspace.dart';
import 'package:fs_frontend/widgets/role_form.dart';
import 'package:fs_frontend/widgets/service_account_details.dart';
import 'package:fs_frontend/widgets/service_account_form.dart';

class WorkspaceSettings extends StatefulWidget {
  const WorkspaceSettings({super.key, required this.workspace});

  final Workspace workspace;

  @override
  State<WorkspaceSettings> createState() => _WorkspaceSettingsState();
}

class _WorkspaceSettingsState extends State<WorkspaceSettings> {
  int selectedIndex = 0;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(
          '${widget.workspace.name} Settings',
          style: const TextStyle(
            fontSize: 24,
            fontWeight: FontWeight.bold,
            color: Color.fromARGB(255, 20, 30, 70),
          ),
        ),
      ),
      body: BlocProvider(
        create: (_) => WSSCubit(
          workspaceRepo: context.read<WorkspaceBloc>().wsRepo,
          workspace: widget.workspace,
        ),
        child: Row(
          children: [
        NavigationRail(
          groupAlignment: 0,
          labelType: NavigationRailLabelType.all,
          onDestinationSelected: (index) => setState(() {
            selectedIndex = index;
          }),
          destinations: const [
            NavigationRailDestination(
              icon: Icon(Icons.person),
              label: Text('Service Accounts'),
            ),
            NavigationRailDestination(
              icon: Icon(Icons.tag),
              label: Text('Roles'),
            ),
            NavigationRailDestination(
              label: Text('Deletion Settings'),
              icon: Icon(Icons.dangerous),
            ),
          ],
          selectedIndex: selectedIndex,
        ),
        Expanded(
          child: Padding(
            padding: const EdgeInsets.all(30),
            child: Builder(
              builder: (_) {
                switch (selectedIndex) {
                  case 0:
                    return WorkspaceServiceAccounts(
                      workspace: widget.workspace,
                    );
                  case 1:
                    return const WorkspaceRoles();
                  case 2:
                    return WorkspaceDangerSettings(
                        workspace: widget.workspace);
                  default:
                    return Container();
                }
              },
            ),
          ),
        )
          ],
        ),
      ),
    );
  }
}

class WorkspaceDangerSettings extends StatefulWidget {
  const WorkspaceDangerSettings({super.key, required this.workspace});
  final Workspace workspace;

  @override
  State<WorkspaceDangerSettings> createState() =>
      _WorkspaceDangerSettingsState();
}

class _WorkspaceDangerSettingsState extends State<WorkspaceDangerSettings> {
  late final TextEditingController workspaceName;

  @override
  void initState() {
    workspaceName = TextEditingController(text: widget.workspace.name);
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          OutlinedButton(
              onPressed: () {
                Navigator.popUntil(
                    context, (route) => route.settings.name == "/ws-select");
                context.read<WorkspaceBloc>().add(
                    WorkspaceDelete(workspaceName: widget.workspace.name));
              },
              child: const Text('Delete Workspace')),
          const SizedBox(height: 10),
          const Text(
              'Deleting workspace will delete all its files and service accounts')
        ],
      ),
    );
  }
}

class WorkspaceServiceAccounts extends StatefulWidget {
  const WorkspaceServiceAccounts({super.key, required this.workspace});
  final Workspace workspace;

  @override
  State<WorkspaceServiceAccounts> createState() =>
      _WorkspaceServiceAccountsState();
}

class _WorkspaceServiceAccountsState extends State<WorkspaceServiceAccounts> {
  @override
  void initState() {
    context.read<WSSCubit>().getServiceAccounts();
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      floatingActionButton: FloatingActionButton(
        onPressed: () => showAddServiceAccountDialog(context),
        child: const Icon(Icons.add),
      ),
      body: BlocBuilder<WSSCubit, WSSState>(
        builder: (context, state) {
          if (state.isLoadingAccounts) {
            return const Center(child: CircularProgressIndicator());
          }
          if (context.runtimeType is WSSServerErr) {
            return const Center(
              child: Text('Server Error'),
            );
          }
          final accounts = state.accounts;
          if (accounts.isEmpty) {
            return const Center(child: Text('No accounts found'));
          }
          return ListView.builder(
            itemCount: accounts.length,
            itemBuilder: (_, index) {
              return Column(
                children: [
                  ListTile(
                    onTap: () => showServiceAccountDetailsDialog(context, accounts[index]),
                    leading: CircleAvatar(child: Text(accounts[index].name[0].toUpperCase())),
                    title: Text(
                      accounts[index].name,
                      style: const TextStyle(fontWeight: FontWeight.bold),
                    ),
                    subtitle: Text(
                        '${accounts[index].username}@${widget.workspace.name}'),
                    trailing: IconButton(
                      onPressed: () {
                        context
                            .read<WSSCubit>()
                            .deleteServiceAccount(accounts[index].username);
                      },
                      icon: const Icon(Icons.delete, size: 20,),
                    ),
                  ),
                  if (index != accounts.length - 1) ...[
                    const Divider(height: 0),
                  ],
                ],
              );
            },
          );
        },
      ),
    );
  }
}

class WorkspaceRoles extends StatefulWidget {
  const WorkspaceRoles({super.key});

  @override
  State<WorkspaceRoles> createState() => _WorkspaceRolesState();
}

class _WorkspaceRolesState extends State<WorkspaceRoles> {
  @override
  void initState() {
    context.read<WSSCubit>().getRoles();
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          showDialog(
              context: context,
              builder: (_) => BlocProvider.value(
                  value: context.read<WSSCubit>(), child: const RoleForm()));
        },
        child: const Icon(Icons.add),
      ),
      body: BlocBuilder<WSSCubit, WSSState>(
        builder: (context, state) {
          if (state.isLoadingRoles) {
            return const Center(child: CircularProgressIndicator());
          }
          if (context.runtimeType is WSSServerErr) {
            return const Center(
              child: Text('Server Error'),
            );
          }
          final roles = state.roles;
          if (roles.isEmpty) {
            return const Center(child: Text('No roles found'));
          }
          return ListView.builder(
            itemCount: roles.length,
            itemBuilder: (_, index) {
              return Column(
                children: [
                  ListTile(
                    trailing: IconButton(
                      onPressed: () {
                        context.read<WSSCubit>().deleteRole(roles[index]);
                      },
                      icon: const Icon(Icons.delete, size: 20,),
                    ),
                    onTap: () {
                      showDialog(
                          context: context,
                          builder: (_) => BlocProvider.value(
                              value: context.read<WSSCubit>(),
                              child: RoleForm(role: roles[index])));
                    },
                    title: Text(roles[index].name),
                    subtitle: roles[index].description.isNotEmpty ? Text(roles[index].description) : null,
                  ),
                  if (index != roles.length - 1) ...[
                    const Divider(height: 0),
                  ],
                ],
              );
            },
          );
        },
      ),
    );
  }
}



