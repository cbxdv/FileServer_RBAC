import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:fs_frontend/blocs/workspace_settings/wss_cubit.dart';
import 'package:fs_frontend/models/account.dart';
import 'package:fs_frontend/models/role.dart';
import 'package:fs_frontend/widgets/add_role_dialog.dart';

showServiceAccountDetailsDialog(BuildContext context, ServiceAccount account) {
  showDialog(context: context, builder: (_) => AlertDialog(
    title: const Text('Update account'),
    content: BlocProvider.value(value: context.read<WSSCubit>(),
        child: ServiceAccountDetails(account: account)),
  ));
}

class ServiceAccountDetails extends StatefulWidget {
  const ServiceAccountDetails({super.key, required this.account});
  final ServiceAccount account;

  @override
  State<ServiceAccountDetails> createState() => _ServiceAccountDetailsState();
}

class _ServiceAccountDetailsState extends State<ServiceAccountDetails> {

  bool isLoading = true;
  late ServiceAccount accountDetails;

  fetchData() async {
    final workspaceName = widget.account.workspace.name;
    final accId = widget.account.id;
    try {
      final roles = await context.read<WSSCubit>().workspaceRepo.getAccountRoles(workspaceName, accId);
      setState(() {
        setState(() {
          isLoading = false;
          accountDetails = widget.account;
          accountDetails.roles = roles;
        });
      });
    } catch (_) {
      setState(() {
        isLoading = false;
      });
    }
  }

  @override
  void initState() {
    setState(() {
      accountDetails = widget.account;
    });
    fetchData();
    super.initState();
  }

  addRoleHandler() async {
    final workspaceName = accountDetails.workspace.name;
    final filterRoles = accountDetails.roles;
    final selectedRole = await showAddRoleDialog(
        context, filterRoles,
        () => context.read<WSSCubit>().workspaceRepo.getRoles(workspaceName)
    );
    if (selectedRole == null) {
      return;
    }
    setState(() {
      isLoading = true;
    });
    try {
      if (context.mounted) {
        await context.read<WSSCubit>().workspaceRepo.assignRole(
            workspaceName, accountDetails.id, selectedRole.id);
      }
      setState(() {
        accountDetails.roles.add(selectedRole);
        isLoading = false;
      });
    } catch (_) {
      setState(() {
        isLoading = true;
      });
    }
  }

  removeRoleHandler(Role role) async {
    final workspaceName = accountDetails.workspace.name;
    final accountId = widget.account.id;
    try {
      await context.read<WSSCubit>().workspaceRepo.unAssignRole(workspaceName, accountId, role.id);
      setState(() {
        widget.account.roles.remove(role);
        isLoading = false;
      });
    } catch (_) {
      setState(() {
        isLoading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    if (isLoading) {
      return const SizedBox(
          height: 100,
          child: Center(child: CircularProgressIndicator())
      );
    }
    return SizedBox(
      width: 400,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisSize: MainAxisSize.min,
        children: [
          const Text('Name', style: TextStyle(fontSize: 13, fontWeight: FontWeight.w700)),
          Text(accountDetails.name, style: const TextStyle(fontSize: 18)),
          const SizedBox(height: 10),
          const Text('Username', style: TextStyle(fontSize: 13, fontWeight: FontWeight.w700)),
          Text(accountDetails.username, style: const TextStyle(fontSize: 18)),
          const SizedBox(height: 10),
          if (accountDetails.roles.isNotEmpty) ...[
            const Text('Roles', style: TextStyle(fontSize: 13, fontWeight: FontWeight.w700)),
            const SizedBox(height: 5),
            Wrap(
              children: List.generate(
                accountDetails.roles.length,
                (index) => Container(
                  margin: const EdgeInsets.symmetric(horizontal: 5),
                  child: Chip(
                    backgroundColor: Colors.transparent,
                    side: BorderSide(color: Colors.grey.shade200),
                    elevation: 0,
                    padding: EdgeInsets.zero,
                    materialTapTargetSize: MaterialTapTargetSize.shrinkWrap,
                    label: Text(accountDetails.roles[index].name),
                    deleteIcon: const Icon(Icons.remove_circle_outline, size: 16,),
                    deleteButtonTooltipMessage: 'Remove',
                    onDeleted: () => removeRoleHandler(accountDetails.roles[index]),
                  )
                ),
              ),
            )
          ],
          const SizedBox(height: 30),
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              OutlinedButton(onPressed: addRoleHandler, child: const Text('Add Role')),
              const SizedBox(width: 10),
              OutlinedButton(onPressed: () {}, child: const Text('Delete')),
              const SizedBox(width: 10),
              OutlinedButton(onPressed: () {
                Navigator.of(context).pop();
              }, child: const Text('Close'))
          ],)
        ],
      ),
    );
  }
}
