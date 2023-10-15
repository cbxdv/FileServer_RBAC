import 'package:flutter/material.dart';
import 'package:fs_frontend/models/role.dart';

Future<Role?> showAddRoleDialog(BuildContext context, List<Role> filterRoles, Function fetchRoles) {
  return showDialog(context: context, builder: (_) {
    return AlertDialog(
      title: const Text('Add Role'),
      content: SizedBox(
        height: 300,
        width: 300,
        child: AddRoleDialog(filterRoles: filterRoles, fetchRoles: fetchRoles),
      ),
      actions: [
        TextButton(onPressed: () {
          Navigator.of(context).pop(null);
        }, child: const Text('Close'))
      ],
      scrollable: true,
    );
  });
}

class AddRoleDialog extends StatefulWidget {
  const AddRoleDialog({super.key, required this.filterRoles, required this.fetchRoles});
  final List<Role> filterRoles;
  final Function fetchRoles;

  @override
  State<AddRoleDialog> createState() => _AddRoleDialogState();
}

class _AddRoleDialogState extends State<AddRoleDialog> {

  bool isLoading = true;
  List<Role> roles = [];
  String err = "";

  void fetchData() async {
    setState(() {
      isLoading = true;
    });
    try {
      final rolesRes = await widget.fetchRoles();
      for (var element in widget.filterRoles) {
        rolesRes.removeWhere((r) => element.id == r.id);
      }
      setState(() {
        roles = rolesRes;
        isLoading = false;
      });
    } catch (_) {
      setState(() {
        err = "Server error";
        isLoading = false;
      });
    }
  }

  void addHandler() async {}

  @override
  void initState() {
    fetchData();
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    if (isLoading) {
      return const Center(child: CircularProgressIndicator(),);
    }
    return ListView.builder(
      itemCount: roles.length,
      itemBuilder: (context, index) {
        return ListTile(title: Text(roles[index].name), onTap: () {
          Navigator.of(context).pop(roles[index]);
        },);
      },
    );
  }
}
