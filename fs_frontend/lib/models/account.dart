// class Account {
//   final String id;
//   final String name;
//   final String username;
//   final bool isOwner;
//
//   Account({
//     required this.id,
//     required this.name,
//     required this.username,
//     required this.isOwner,
//   });
//
//   factory Account.fromJson(Map<String, dynamic> json) {
//     return Account(
//       id: json['id'],
//       name: json['name'],
//       isOwner: json['isOwner'],
//       username: json['username'],
//     );
//   }
// }

import 'package:fs_frontend/models/role.dart';
import 'package:fs_frontend/models/workspace.dart';

abstract class Account {
  final String id;
  final String name;
  final String username;
  Account({required this.id, required this.name, required this.username});
}

class OwnerAccount extends Account {
  OwnerAccount(
      {required super.id, required super.name, required super.username});

  factory OwnerAccount.fromJson(Map<String, dynamic> json) {
    return OwnerAccount(
      id: json['id'],
      name: json['name'],
      username: json['username'],
    );
  }
}

class ServiceAccount extends Account {
  final Workspace workspace;
  List<Role> roles = [];

  ServiceAccount(
      {required this.workspace,
      required super.id,
      required super.name,
      required super.username});

  factory ServiceAccount.fromJson(Map<String, dynamic> json) {
    return ServiceAccount(
      id: json['id'],
      name: json['name'],
      username: json['username'],
      workspace: Workspace(name: json['workspace']),
    );
  }
}
