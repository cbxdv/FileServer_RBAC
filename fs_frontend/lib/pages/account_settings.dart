import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:fs_frontend/blocs/auth/auth_bloc.dart';
import 'package:fs_frontend/widgets/styled_text_field.dart';

class AccountSettings extends StatefulWidget {
  const AccountSettings({super.key});

  @override
  State<AccountSettings> createState() => _AccountSettingsState();
}

class _AccountSettingsState extends State<AccountSettings> {
  late final TextEditingController oldPassword;
  late final TextEditingController newPassword;
  late final TextEditingController confirmNewPassword;

  @override
  void initState() {
    oldPassword = TextEditingController();
    newPassword = TextEditingController();
    confirmNewPassword = TextEditingController();

    super.initState();
  }

  changePasswordHandler() async {
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (context) => const Center(child: CircularProgressIndicator(color: Colors.white)),
    );
    await context.read<AuthBloc>().authRepo.changePassword(
        oldPassword: oldPassword.text, newPassword: newPassword.text);
    if (context.mounted) {
      Navigator.of(context).pop();
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Account Settings',
            style: TextStyle(
              fontSize: 24,
              fontWeight: FontWeight.bold,
              color: Color.fromARGB(255, 20, 30, 70),
            )),
      ),
      body: Padding(
        padding: const EdgeInsets.all(30.0),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          crossAxisAlignment: CrossAxisAlignment.center,
          children: [
            const Text(
              'Change Password',
              style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold, color: Color.fromARGB(255, 20, 30, 70),),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 30),
            Center(
              child: SizedBox(
                width: 500,
                child: StyledTextField(
                    name: 'Old Password',
                    controller: oldPassword,
                    isPassword: true),
              ),
            ),
            Center(
              child: SizedBox(
                width: 500,
                child: StyledTextField(
                  name: 'New Password',
                  controller: newPassword,
                  isPassword: true,
                ),
              ),
            ),
            Center(
              child: SizedBox(
                width: 500,
                child: StyledTextField(
                  name: 'Confirm New Password',
                  controller: confirmNewPassword,
                  isPassword: true,
                ),
              ),
            ),
            const SizedBox(height: 30),
            Center(
              child: OutlinedButton(
                onPressed: changePasswordHandler,
                child: const Text('Change Password'),
              ),
            ),
            // SizedBox(
            //   height: 30,
            // ),
            // Center(
            //   child: OutlinedButton(
            //     onPressed: () {},
            //     child: Text('Delete Account'),
            //   ),
            // )
          ],
        ),
      ),
    );
  }
}
