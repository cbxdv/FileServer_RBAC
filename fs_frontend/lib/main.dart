import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:fs_frontend/blocs/auth/auth_bloc.dart';
import 'package:fs_frontend/blocs/transfer/transfer_bloc.dart';
import 'package:fs_frontend/blocs/workspace/workspace_bloc.dart';
import 'package:fs_frontend/models/account.dart';
import 'package:fs_frontend/pages/explorer.dart';
import 'package:fs_frontend/pages/login.dart';
import 'package:fs_frontend/pages/splash.dart';
import 'package:fs_frontend/pages/workspace_select.dart';
import 'package:fs_frontend/repos/transfer_repo.dart';

void main() {
  WidgetsFlutterBinding.ensureInitialized();
  runApp(BlocProvider(
    create: (_) => AuthBloc(),
    child: const MyApp(),
  ));
}

class MyApp extends StatefulWidget {
  const MyApp({super.key});

  @override
  State<MyApp> createState() => _MyAppState();
}

class _MyAppState extends State<MyApp> {
  final _navigatorKey = GlobalKey<NavigatorState>();

  NavigatorState get _navigator => _navigatorKey.currentState!;

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      navigatorKey: _navigatorKey,
      title: 'FS',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(
          seedColor: const Color.fromARGB(255, 244, 238, 238),
        ),
        useMaterial3: true,
      ),
      builder: (context, child) {
        return BlocListener<AuthBloc, AuthState>(
          listener: (c, state) {
            if (state is AuthLoggedInState) {
              if (state.account is OwnerAccount) {
                _navigator.pushAndRemoveUntil(
                  MaterialPageRoute(
                      builder: (_) => BlocProvider(
                            create: (_) => WorkspaceBloc(
                                secureStorage:
                                    context.read<AuthBloc>().secureStorage),
                            child: const WorkspaceSelect(),
                          ),
                      settings: const RouteSettings(name: "/ws-select")),
                  (route) => false,
                );
              } else if (state.account is ServiceAccount) {
                _navigator.pushAndRemoveUntil(
                  MaterialPageRoute(
                      builder: (_) => BlocProvider(
                            create: (_) => TransferBloc(
                              transferRepo: TransferRepo(
                                secureStorage:
                                    context.read<AuthBloc>().secureStorage,
                              ),
                            ),
                            child: Explorer(
                                currentWorkspace:
                                    (state.account as ServiceAccount)
                                        .workspace),
                          ),
                      settings: const RouteSettings(name: "/ws-select")),
                  (route) => false,
                );
              }
            }
            if (state is AuthLoggedOutState) {
              _navigator.pushAndRemoveUntil(
                MaterialPageRoute(builder: (_) => const Login()),
                (route) => false,
              );
            }
          },
          child: child,
        );
      },
      onGenerateRoute: (_) => MaterialPageRoute(builder: (_) => const Splash()),
    );
  }
}
