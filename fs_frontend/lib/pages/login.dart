import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:fs_frontend/blocs/auth/auth_bloc.dart';
import 'package:fs_frontend/pages/register.dart';
import 'package:fs_frontend/widgets/styled_text_field.dart';

class Login extends StatefulWidget {
  const Login({super.key});

  @override
  State<Login> createState() => _LoginState();
}

class _LoginState extends State<Login> {
  @override
  Widget build(BuildContext context) {
    if (context.read<AuthBloc>().state is AuthLoggedInState) {
      Navigator.of(context).pop();
    }
    return Scaffold(
      body: DefaultTabController(
        length: 2,
        child: BlocListener(
          bloc: context.read<AuthBloc>(),
          listener: (context, state) {
            if (state is AuthLoginFailureState) {
              showDialog(context: context, builder: (context) => AlertDialog(
                 title: Text(state.title),
                content: Text(state.description),
                actions: [
                  TextButton(onPressed: () {Navigator.pop(context);}, child: const Text('OK'))
                ],
              ));
            }
          },
          child: const Center(
            child: SizedBox(
              width: 400,
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Row(
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
                  SizedBox(height: 60),
                  TabBar(
                    tabs: [
                      Tab(child: Text('Login as Owner')),
                      Tab(child: Text('Login as Service'))
                    ],
                    labelColor: Color.fromARGB(255, 20, 30, 70),
                    indicatorColor: Color.fromARGB(255, 20, 30, 70),
                    unselectedLabelColor: Color.fromARGB(255, 159, 159, 159),
                  ),
                  SizedBox(height: 30),
                  SizedBox(
                    height: 300,
                    width: 400,
                    child: TabBarView(
                      children: [
                        OwnerAccountTab(),
                        ServiceAccountTab(),
                      ],
                    ),
                  )
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class OwnerAccountTab extends StatefulWidget {
  const OwnerAccountTab({super.key});

  @override
  State<OwnerAccountTab> createState() => _OwnerAccountTabState();
}

class _OwnerAccountTabState extends State<OwnerAccountTab> {
  late TextEditingController email;
  late TextEditingController password;

  @override
  void initState() {
    email = TextEditingController();
    password = TextEditingController();
    super.initState();
  }

  @override
  void dispose() {
    email.dispose();
    password.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 10),
      child: Column(
        children: [
          StyledTextField(
              name: 'Email',
              controller: email,
              inputType: TextInputType.emailAddress),
          StyledTextField(
              name: 'Password', controller: password, isPassword: true),
          Container(
            margin: const EdgeInsets.symmetric(vertical: 30),
            child: SizedBox(
              height: 40,
              width: 200,
              child: TextButton(
                onPressed: () {
                  context.read<AuthBloc>().add(OwnerAccountLoginEvent(
                        email: email.text,
                        password: password.text,
                      ));
                },
                style: TextButton.styleFrom(
                  backgroundColor: const Color.fromARGB(255, 20, 30, 70),
                ),
                child: BlocBuilder<AuthBloc, AuthState>(
                  builder: (BuildContext context, AuthState state) {
                    if (state is AuthLoggingInState) {
                      return const SizedBox(
                          height: 10,
                          width: 10,
                          child: CircularProgressIndicator(
                            color: Colors.white,
                            strokeWidth: 2,
                          ));
                    }
                    return const Text(
                      'Login',
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
              const Text('Not having an account?'),
              TextButton(
                onPressed: () {
                  Navigator.of(context).push(
                    MaterialPageRoute(builder: (_) => const Register())
                  );
                },
                child: const Text('Register'),
              )
            ],
          )
        ],
      ),
    );
  }
}

class ServiceAccountTab extends StatefulWidget {
  const ServiceAccountTab({super.key});

  @override
  State<ServiceAccountTab> createState() => _ServiceAccountTabState();
}

class _ServiceAccountTabState extends State<ServiceAccountTab> {
  late TextEditingController username;
  late TextEditingController password;

  @override
  void initState() {
    username = TextEditingController();
    password = TextEditingController();
    super.initState();
  }

  @override
  void dispose() {
    username.dispose();
    password.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 10),
      child: Column(
        children: [
          StyledTextField(name: 'Username', controller: username),
          StyledTextField(
              name: 'Password', controller: password, isPassword: true),
          Container(
            margin: const EdgeInsets.symmetric(vertical: 30),
            child: SizedBox(
              height: 40,
              width: 200,
              child: TextButton(
                onPressed: () {
                  context.read<AuthBloc>().add(
                        ServiceAccountLoginEvent(
                          username: username.text,
                          password: password.text,
                        ),
                      );
                },
                style: TextButton.styleFrom(
                  backgroundColor: const Color.fromARGB(255, 20, 30, 70),
                ),
                child: BlocBuilder<AuthBloc, AuthState>(
                  builder: (BuildContext context, AuthState state) {
                    if (state is AuthLoggingInState) {
                      return const SizedBox(
                          height: 10,
                          width: 10,
                          child: CircularProgressIndicator(
                            color: Colors.white,
                            strokeWidth: 2,
                          ));
                    }
                    return const Text(
                      'Login',
                      style: TextStyle(color: Colors.white),
                    );
                  },
                ),
              ),
            ),
          )
        ],
      ),
    );
  }
}
