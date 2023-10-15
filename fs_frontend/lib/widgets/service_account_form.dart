import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:fs_frontend/blocs/workspace_settings/wss_cubit.dart';
import 'package:fs_frontend/widgets/styled_text_field.dart';

showAddServiceAccountDialog(BuildContext context) {
  final TextEditingController name = TextEditingController();
  final TextEditingController username = TextEditingController();
  final TextEditingController password = TextEditingController();

  showDialog(
    context: context,
    builder: (_) => AlertDialog(
      title: const Text('Add Service Account'),
      content: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          SizedBox(
              width: 500,
              child: StyledTextField(name: 'Name', controller: name)),
          SizedBox(
              width: 500,
              child: StyledTextField(name: 'Username', controller: username)),
          SizedBox(
              width: 500,
              child: StyledTextField(
                name: 'Password',
                controller: password,
                isPassword: true,
              )),
        ],
      ),
      actions: [
        TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('Discard')),
        TextButton(
            onPressed: () {
              Navigator.pop(context);
              context.read<WSSCubit>().createServiceAccount(
                  name.text, username.text, password.text);
            },
            child: const Text('Create')),
      ],
    ),
  );
}

