class ApiConstants {
  static String host = "http://localhost:3000";

  static String registerOwnerAccount = "$host/auth/register";
  static String loginOwnerAccount = "$host/auth/login";
  static String checkAuth = "$host/auth/check";
  static String changePassword = "$host/auth/change-password";

  static String workspacesOperations = "$host/ws/op";

  static String dirQuery = "$host/fs/dir/query";
  static String fileQuery = "$host/fs/file/query";

  static String dirDetails = "$host/fs/dir/details";
  static String fileDetails = "$host/fs/file/details";

  static String workspaceAccountOperations = "$host/ws/account";

  static String allRolesQuery = "$host/roles/details";
  static String roleOperations = "$host/role/op";

  static String upload = "$host/fs/upload";
  static String download = "$host/fs/download";

  static String loginServiceAccount = "$host/auth/sa/login";

  static String sharedContent = "$host/fs/shared/query";

  static String assignRoleFS = "$host/rbac/fs";

  static String accountRoles = "$host/roles/sa";
  static String accountRoleAssign = "$host/role/assign";
}