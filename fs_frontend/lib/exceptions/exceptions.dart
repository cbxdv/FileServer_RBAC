class AccountAlreadyExists implements Exception {
  final String description;
  AccountAlreadyExists(this.description);
}

class WeakPassword implements Exception {
  final String description;
  WeakPassword(this.description);
}

class InvalidCredentials implements Exception {
  final String description;
  InvalidCredentials(this.description);
}

class ServerError implements Exception {
  String description = "Server error";
}

class UnauthorizedError implements Exception {}

class WorkspaceAlreadyExists implements Exception {}

class ResourceAlreadyExists implements Exception {
  final String description;
  ResourceAlreadyExists(this.description);
}

class ChildrenExists implements Exception {
  final String description;
  ChildrenExists({required this.description});
}

class PermissionDenied implements Exception {}