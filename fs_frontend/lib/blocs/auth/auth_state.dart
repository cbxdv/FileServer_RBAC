part of 'auth_bloc.dart';

abstract class AuthState {}

abstract class AuthActionState extends AuthState {}

class AuthInitialState extends AuthState {}

class AuthLoadingState extends AuthState {}

class AuthLoggingInState extends AuthState {}

class AuthRegisteringState extends AuthState {}

class AuthLoggedInState extends AuthState {
  final Account account;

  AuthLoggedInState({required this.account});
}

class AuthRegisterErrorState extends AuthState {
  final String title;
  final String description;
  AuthRegisterErrorState({required this.title, required this.description});
}

class AuthServerErrorState extends AuthState {}

class AuthLoginFailureState extends AuthActionState {
  final String title;
  final String description;
  AuthLoginFailureState({required this.title, required this.description});
}

class AuthLoggedOutState extends AuthState {}
