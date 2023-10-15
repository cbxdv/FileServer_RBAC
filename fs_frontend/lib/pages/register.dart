import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:fs_frontend/blocs/auth/auth_bloc.dart';
import 'package:fs_frontend/pages/login.dart';
import 'package:fs_frontend/widgets/styled_text_field.dart';

class Register extends StatefulWidget {
  const Register({super.key});

  @override
  State<Register> createState() => _RegisterState();
}

class _RegisterState extends State<Register> {
  late TextEditingController nameController;
  late TextEditingController emailController;
  late TextEditingController passwordController;
  late TextEditingController confirmPasswordController;

  @override
  void initState() {
    nameController = TextEditingController();
    emailController = TextEditingController();
    passwordController = TextEditingController();
    confirmPasswordController = TextEditingController();
    super.initState();
  }

  @override
  void dispose() {
    nameController.dispose();
    emailController.dispose();
    passwordController.dispose();
    confirmPasswordController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    if (context.read<AuthBloc>().state is AuthLoggedInState) {
      Navigator.of(context).pop();
    }
    return Scaffold(
      body: BlocListener(
        bloc: context.read<AuthBloc>(),
        listener: (BuildContext context, state) {
          if (state is AuthServerErrorState) {
            showDialog(context: context, builder: (context) => AlertDialog(
              title: const Text("Auth Server Error"),
              content: const Text("Error registering account. Try again later"),
              actions: [
                TextButton(onPressed: () {Navigator.pop(context);}, child: const Text('OK'))
              ],
            ));
          }
          if (state is AuthRegisterErrorState) {
            showDialog(context: context, builder: (context) => AlertDialog(
              title: Text(state.title),
              content: Text(state.description),
              actions: [
                TextButton(onPressed: () {Navigator.pop(context);}, child: const Text('OK'))
              ],
            ));
          }
        },
        child: Center(
          child: SingleChildScrollView(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                const Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Icon(
                      Icons.cloud_circle,
                      size: 50,
                      color: Color.fromARGB(255, 20, 30, 70),
                    ),
                    SizedBox(width: 10),
                    Text(
                      'File Server - RBAC',
                      style: TextStyle(
                        fontWeight: FontWeight.bold,
                        fontSize: 16,
                        color: Color.fromARGB(255, 20, 30, 70),
                      ),
                    )
                  ],
                ),
                const SizedBox(height: 60),
                Center(
                  child: SizedBox(
                    width: 400,
                    child: Column(
                      children: [
                        StyledTextField(
                          name: 'Name',
                          controller: nameController,
                        ),
                        StyledTextField(
                          name: 'Email',
                          controller: emailController,
                          inputType: TextInputType.emailAddress,
                        ),
                        StyledTextField(
                          name: 'Password',
                          controller: passwordController,
                          isPassword: true,
                        ),
                        StyledTextField(
                          name: 'Password Again',
                          controller: confirmPasswordController,
                          isPassword: true,
                        ),
                        Container(
                          margin: const EdgeInsets.symmetric(vertical: 30),
                          child: SizedBox(
                            height: 40,
                            width: 200,
                            child: TextButton(
                              onPressed: () {
                                context.read<AuthBloc>().add(
                                      RegisterAccountEvent(
                                        name: nameController.text,
                                        email: emailController.text,
                                        password: passwordController.text,
                                      ),
                                    );
                              },
                              style: TextButton.styleFrom(
                                backgroundColor: const Color.fromARGB(255, 20, 30, 70),
                              ),
                              child: BlocBuilder<AuthBloc, AuthState>(
                                builder: (BuildContext context, AuthState state) {
                                  if (state is AuthRegisteringState) {
                                    return const SizedBox(
                                      height: 10,
                                      width: 10,
                                      child: CircularProgressIndicator(
                                        color: Colors.white,
                                        strokeWidth: 2,
                                      ),
                                    );
                                  }
                                  return const Text(
                                    'Register',
                                    style: TextStyle(color: Colors.white),
                                  );
                                },
                              ),
                            ),
                          ),
                        ),
                        Row(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            const Text('Already having an account?'),
                            TextButton(
                              onPressed: () {
                                Navigator.of(context).push(
                                    MaterialPageRoute(builder: (_) => const Login())
                                );
                              },
                              child: const Text('Login'),
                            )
                          ],
                        )
                      ],
                    ),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
