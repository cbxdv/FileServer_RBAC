import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:fs_frontend/blocs/auth/auth_bloc.dart';
import 'package:fs_frontend/blocs/fs/fs_bloc.dart';
import 'package:fs_frontend/blocs/workspace/workspace_bloc.dart';
import 'package:fs_frontend/models/account.dart';
import 'package:fs_frontend/pages/account_settings.dart';
import 'package:fs_frontend/pages/workspace_settings.dart';
import 'package:fs_frontend/widgets/transfer_pane.dart';

class ExplorerSidebar extends StatelessWidget {
  const ExplorerSidebar({super.key});

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
              child: BlocBuilder<FSBloc, FSState>(
                builder: (_, fsState) => Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    if ((context.read<AuthBloc>().state as AuthLoggedInState).account is OwnerAccount) ...[
                      IconButton(
                        onPressed: () {
                          Navigator.of(context).pop();
                        },
                        icon: const Icon(Icons.workspaces_outline),
                        tooltip: 'Workspaces',
                        iconSize: 30,
                      )
                    ],
                    const SizedBox(height: 30),
                    IconButton(
                      onPressed: () {
                        showTransferPane(context);
                      },
                      icon: const Icon(Icons.downloading),
                      tooltip: 'Transfers',
                      iconSize: 30,
                    ),
                    if ((context.read<AuthBloc>().state as AuthLoggedInState).account is OwnerAccount) ...[
    const SizedBox(height: 30),
    IconButton(
    onPressed: () {
    Navigator.of(context).push(MaterialPageRoute(
    builder: (_) => BlocProvider.value(
    value: context.read<WorkspaceBloc>(),
    child: WorkspaceSettings(
    workspace: context.read<FSBloc>().workspace),
    ),
    ));
    },
    icon: const Icon(Icons.admin_panel_settings_outlined),
    tooltip: 'Workspace Settings',
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
    )
                    ],
                  ],
                ),
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