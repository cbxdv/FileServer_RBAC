part of 'auth_bloc.dart';

abstract class AuthEvent {}

class CheckAuthEvent extends AuthEvent {}

class OwnerAccountLoginEvent extends AuthEvent {
  final String email;
  final String password;

  OwnerAccountLoginEvent({required this.email, required this.password});
}

class ServiceAccountLoginEvent extends AuthEvent {
  final String username;
  final String password;

  ServiceAccountLoginEvent({required this.username, required this.password});
}

class LogoutEvent extends AuthEvent {}

class RegisterAccountEvent extends AuthEvent {
  final String name;
  final String email;
  final String password;

  RegisterAccountEvent({required this.name, required this.email, required this.password});
}
