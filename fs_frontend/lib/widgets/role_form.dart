import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:fs_frontend/blocs/workspace_settings/wss_cubit.dart';
import 'package:fs_frontend/models/role.dart';
import 'package:fs_frontend/widgets/styled_text_field.dart';

class RoleForm extends StatefulWidget {
  const RoleForm({super.key, this.role});
  final Role? role;

  @override
  State<RoleForm> createState() => _RoleFormState();
}

class _RoleFormState extends State<RoleForm> {
  late final TextEditingController name;
  late final TextEditingController description;
  bool canRead = false;
  bool canCreate = false;
  bool canDelete = false;

  @override
  void initState() {
    name = TextEditingController();
    description = TextEditingController();

    if (widget.role != null) {
      name.text = widget.role?.name ?? "";
      description.text = widget.role?.description ?? "";
      canRead = widget.role?.canRead ?? false;
      canCreate = widget.role?.canCreate ?? false;
      canDelete = widget.role?.canDelete ?? false;
    }

    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: Text(widget.role == null ? 'Add Role' : 'Update Role'),
      content: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          SizedBox(
              width: 500,
              child: StyledTextField(name: 'Name', controller: name)),
          SizedBox(
              width: 500,
              child: StyledTextField(
                  name: 'Description', controller: description)),
          const SizedBox(height: 10),
          Wrap(
            spacing: 5,
            children: [
              ChoiceChip(
                  label: const Text('Read'),
                  selected: canRead,
                  onSelected: (value) => setState(() {
                    canRead = value;
                  })),
              ChoiceChip(
                  label: const Text('Create'),
                  selected: canCreate,
                  onSelected: (value) => setState(() {
                    canCreate = value;
                  })),
              ChoiceChip(
                  label: const Text('Delete'),
                  selected: canDelete,
                  onSelected: (value) => setState(() {
                    canDelete = value;
                  }))
            ],
          )
        ],
      ),
      actions: [
        TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('Discard')),
        TextButton(
            onPressed: () {
              if (widget.role == null) {
                context.read<WSSCubit>().createRole(
                  name: name.text,
                  description: description.text,
                  canRead: canRead,
                  canCreate: canCreate,
                  canDelete: canDelete,
                );
              } else {
                context.read<WSSCubit>().updateRole(
                  id: widget.role?.id ?? "",
                  name: name.text,
                  description: description.text,
                  canRead: canRead,
                  canCreate: canCreate,
                  canDelete: canDelete,
                );
              }
              Navigator.of(context).pop();
            },
            child: Text(widget.role == null ? 'Create' : 'Update')),
      ],
    );
  }
}
