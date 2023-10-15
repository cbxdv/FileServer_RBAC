import 'package:flutter/material.dart';

class StyledTextField extends StatefulWidget {
  const StyledTextField({
    super.key,
    required this.name,
    required this.controller,
    this.inputType = TextInputType.text,
    this.isPassword = false,
  });

  final String name;
  final TextEditingController controller;
  final TextInputType inputType;
  final bool isPassword;

  @override
  State<StyledTextField> createState() => _StyledTextFieldState();
}

class _StyledTextFieldState extends State<StyledTextField> {
  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.symmetric(vertical: 10),
      child: TextField(
        controller: widget.controller,
        keyboardType: TextInputType.emailAddress,
        obscureText: widget.isPassword,
        autocorrect: !widget.isPassword,
        enableSuggestions: !widget.isPassword,
        decoration: InputDecoration(
          hintText: widget.name,
          filled: true,
          fillColor: const Color.fromARGB(255, 244, 238, 238),
          border: OutlineInputBorder(
            borderSide: BorderSide.none,
            borderRadius: BorderRadius.circular(20),
          ),
        ),
        style: const TextStyle(
          fontSize: 14
        ),
      ),
    );
  }
}
