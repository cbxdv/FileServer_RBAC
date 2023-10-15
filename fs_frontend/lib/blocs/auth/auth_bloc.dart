import 'dart:async';

import 'package:bloc/bloc.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:fs_frontend/exceptions/exceptions.dart';

import 'package:fs_frontend/models/account.dart';
import 'package:fs_frontend/repos/auth_repo.dart';

part 'auth_event.dart';

part 'auth_state.dart';

class AuthBloc extends Bloc<AuthEvent, AuthState> {
  late final AuthRepo authRepo;
  late final FlutterSecureStorage secureStorage;

  AuthBloc() : super(AuthInitialState()) {
    secureStorage = const FlutterSecureStorage();
    authRepo = AuthRepo(secureStorage: secureStorage);

    on<CheckAuthEvent>(checkAuth);
    on<OwnerAccountLoginEvent>(ownerAccountLogin);
    on<ServiceAccountLoginEvent>(serviceAccountLogin);
    on<RegisterAccountEvent>(register);
    on<LogoutEvent>(logout);
  }

  FutureOr<void> checkAuth(
    CheckAuthEvent event,
    Emitter<AuthState> emit,
  ) async {
    emit(AuthLoadingState());
    try {
      Account? account = await authRepo.checkAuth();
      if (account == null) {
        emit(AuthLoggedOutState());
      } else {
        emit(AuthLoggedInState(account: account));
      }
    } catch (_) {
      emit(AuthLoggedOutState());
    }
  }

  FutureOr<void> ownerAccountLogin(
    OwnerAccountLoginEvent event,
    Emitter<AuthState> emit,
  ) async {
    if (state is AuthLoggingInState) {
      return;
    }
    emit(AuthLoggingInState());
    if (event.email.isEmpty || event.password.isEmpty) {
      emit(
        AuthLoginFailureState(
            title: "Invalid credentials",
            description:
                "The provided username or password is invalid. Check again."),
      );
      return;
    }
    try {
      OwnerAccount? account = await authRepo.ownerAccountLogin(
        email: event.email,
        password: event.password,
      );
      if (account == null) {
        emit(AuthInitialState());
      } else {
        emit(AuthLoggedInState(account: account));
      }
    } on InvalidCredentials {
      emit(
        AuthLoginFailureState(
            title: "Invalid credentials",
            description:
                "The provided username or password is invalid. Check again."),
      );
    } catch (e) {
      emit(AuthLoginFailureState(
        title: "Error",
        description: "Auth server error. Try again later"
      ));
    }
  }

  FutureOr<void> serviceAccountLogin(
    ServiceAccountLoginEvent event,
    Emitter<AuthState> emit,
  ) async {
    if (state is AuthLoggingInState) {
      return;
    }
    emit(AuthLoggingInState());
    if (event.username.isEmpty || event.password.isEmpty) {
      emit(
        AuthLoginFailureState(
            title: "Invalid credentials",
            description:
            "The provided username or password is invalid. Check again."),
      );
      return;
    }
    try {
      ServiceAccount? account = await authRepo.serviceAccountLogin(
        username: event.username,
        password: event.password,
      );
      if (account == null) {
        emit(AuthInitialState());
      } else {
        emit(AuthLoggedInState(account: account));
      }
    } on InvalidCredentials {
      emit(
        AuthLoginFailureState(
            title: "Invalid credentials",
            description:
            "The provided username or password is invalid. Check again."),
      );
    } catch (e) {
      emit(AuthLoginFailureState(
          title: "Error",
          description: "Auth server error. Try again later"
      ));
    }
  }

  FutureOr<void> logout(
    LogoutEvent event,
    Emitter<AuthState> emit,
  ) async {
    emit(AuthLoadingState());
    authRepo.logout();
    emit(AuthLoggedOutState());
  }

  FutureOr<void> register(
    RegisterAccountEvent event,
    Emitter<AuthState> emit,
  ) async {
    if (event.name.isEmpty || event.email.isEmpty || event.password.isEmpty) {
      return;
    }
    if (state is AuthRegisteringState) {
      return;
    }
    emit(AuthRegisteringState());
    try {
      Account? account = await authRepo.registerAccount(
        name: event.name,
        email: event.email,
        password: event.password,
      );
      if (account == null) {
        emit(AuthInitialState());
      } else {
        emit(AuthLoggedInState(account: account));
      }
    } on WeakPassword catch (e) {
      emit(AuthRegisterErrorState(
          title: 'Weak Password', description: e.description));
    } on AccountAlreadyExists catch (e) {
      emit(AuthRegisterErrorState(
          title: 'Account Already Exists', description: e.description));
    } catch(_) {
      emit(AuthServerErrorState());
    }
  }
}
