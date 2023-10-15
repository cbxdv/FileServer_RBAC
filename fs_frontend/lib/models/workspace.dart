class Workspace {
  final String name;

  Workspace({required this.name});

  factory Workspace.fromJson(Map<String, dynamic> json) {
    return Workspace(name: json['name']);
  }
}