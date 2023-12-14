// ignore_for_file: must_be_immutable

import 'package:equatable/equatable.dart';
import 'package:fs_frontend/models/account.dart';
import 'package:fs_frontend/models/role.dart';

class WSSState extends Equatable {
  bool isUpdatingWorkspace = false;
  bool isLoadingAccounts = false;
  bool isLoadingRoles = false;
  List<ServiceAccount> accounts = [];
  List<Role> roles = [];

  WSSState();

  @override
  List<Object?> get props {
    List<Object> props = [
      isUpdatingWorkspace,
      isLoadingRoles,
      isLoadingAccounts,
      accounts,
      roles,
    ];
    props.addAll(accounts);
    props.addAll(roles);
    return props;
  }

  factory WSSState.copyFrom(WSSState old) {
    final wss = WSSState();
    wss.isUpdatingWorkspace = old.isUpdatingWorkspace;
    wss.isLoadingAccounts = old.isLoadingAccounts;
    wss.accounts = old.accounts;
    wss.roles = old.roles;
    return wss;
  }
}
