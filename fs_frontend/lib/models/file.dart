class FileModel {
  final String id;
  final String name;
  final int size;
  final String location;
  final DateTime createdOn;

  FileModel({
    required this.id,
    required this.name,
    required this.size,
    required this.location,
    required this.createdOn,
  });

  factory FileModel.fromJson(Map<String, dynamic> json) {
    return FileModel(
      id: json['id'],
        name: json['name'],
        size: int.parse(json['size'].toString()),
        location: json['location'],
        createdOn: DateTime.parse(json['createdOn']),
    );
  }
}
