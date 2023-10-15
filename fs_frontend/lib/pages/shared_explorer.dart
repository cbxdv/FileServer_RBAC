import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:fs_frontend/blocs/fs/fs_bloc.dart';
import 'package:fs_frontend/models/directory.dart';
import 'package:fs_frontend/models/file.dart';
import 'package:fs_frontend/widgets/directory_row.dart';
import 'package:fs_frontend/widgets/explorer_sidebar.dart';
import 'package:fs_frontend/widgets/file_row.dart';

class SharedExplorer extends StatefulWidget {
  const SharedExplorer({super.key, required this.workspaceName});
  final String workspaceName;

  @override
  State<SharedExplorer> createState() => _SharedExplorerState();
}

class _SharedExplorerState extends State<SharedExplorer> {
  bool isLoading = false;
  List<dynamic> contents = [];

  fetchData() async {
    try {
      setState(() {
        isLoading = true;
      });
      final data = await context.read<FSBloc>().fsRepo.getShared(widget.workspaceName);
      setState(() {
        contents = data;
        isLoading = false;
      });
    } catch (_) {
      setState(() {
        isLoading = false;
      });
    }
  }

  @override
  void initState() {
    fetchData();
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Row(
        children: [
          const ExplorerSidebar(),
          Expanded(
            child: Padding(
              padding: const EdgeInsets.all(30.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Padding(
                    padding: EdgeInsets.symmetric(horizontal: 30.0),
                    child: Text("Shared", style: TextStyle(
                      fontWeight: FontWeight.bold,
                      fontSize: 20,
                      color: Color.fromARGB(255, 20, 30, 70),
                    ),),
                  ),
                  const SizedBox(height: 30),
                  Expanded(
                    child: Builder(builder: (context) {
                      if (isLoading) {
                        return const Center(child: CircularProgressIndicator());
                      }
                      if (contents.isEmpty) {
                        return const Center(child: Text('No data found'),);
                      }
                      return ListView.builder(
                        itemCount: contents.length,
                        itemBuilder: (context, index) {
                          if (contents[index] is DirectoryModel) {
                            return DirectoryRow(
                              directory: contents[index],
                            );
                          } else if (contents[index] is FileModel) {
                            return FileRow(file: contents[index]);
                          }
                          return Container();
                        },
                      );
                    })
                  )
                ],
              ),
            ),
          ),
        ]
      )
    );
  }
}
