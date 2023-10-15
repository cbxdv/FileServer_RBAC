import 'package:fs_frontend/models/file.dart';
import 'package:fs_frontend/models/role.dart';

class DirectoryModel {
  final String id;
  final String name;
  final String location;
  final DateTime createdOn;
  List<dynamic> contents = [];
  List<Role> roles = [];

  DirectoryModel({
    required this.id,
    required this.name,
    required this.location,
    required this.createdOn,
  });

  factory DirectoryModel.fromJson(Map<String, dynamic> json) {
    DirectoryModel newDir = DirectoryModel(
      id: json['id'],
      name: json['name'],
      location: json['location'],
      createdOn: DateTime.parse(json['createdOn']),
    );
    if (json['contents'] == null) {
      return newDir;
    }
    final jsonList = json['contents'] as List;
    for (var element in jsonList) {
      if (element['type'] == 'file') {
        newDir.contents.add(FileModel.fromJson(element));
      }
      else if (element['type'] == 'directory') {
        newDir.contents.add(DirectoryModel.fromJson(element));
      }
    }
    return newDir;
  }
}
