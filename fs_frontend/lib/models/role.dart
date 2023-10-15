class Role {
  final String id;
  final String name;
  final String description;
  final bool canRead;
  final bool canCreate;
  final bool canDelete;

  Role({
    required this.id,
    required this.name,
    required this.description,
    required this.canRead,
    required this.canCreate,
    required this.canDelete,
  });

  factory Role.fromJson(Map<String, dynamic> json) {
    return Role(
      id: json['id'],
      name: json['name'],
      description: json['description'],
      canRead: json['canRead'] as bool,
      canCreate: json['canCreate'] as bool,
      canDelete: json['canDelete'] as bool,
    );
  }
}
