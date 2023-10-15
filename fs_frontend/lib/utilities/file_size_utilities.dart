getSizeString(int size) {
  double s = size / 1024;
  // In KB
  if (s <= 1024) {
    return '${s.toStringAsFixed(2)} KB';
  }
  // In MB
  s = size / 1025 / 1024;
  if (s <= 1024) {
    return '${s.toStringAsFixed(2)} MB';
  }
  // In GB
  s = size / 1024 / 1024 / 1024;
  return '${s.toStringAsFixed(2)} GB';
}