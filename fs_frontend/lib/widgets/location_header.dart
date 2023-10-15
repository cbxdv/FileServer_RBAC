import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:fs_frontend/blocs/fs/fs_bloc.dart';
import 'package:fs_frontend/widgets/directory_row.dart';

class LocationHeader extends StatelessWidget {
  const LocationHeader({super.key, required this.location});
  final String location;

  @override
  Widget build(BuildContext context) {

    showFolderDetails(location) async {
      final dirDetails = await context.read<FSBloc>().fsRepo.getDirectoryDetails(location);
      if (context.mounted) {
        showDirectoryDialog(context, dirDetails);

      }
    }

    List<Widget> getFolderPath() {
      List<Widget> widgets = [];
      List<String> locationSplit = location.split('/');
      for (var i = 0; i < locationSplit.length; i++) {
        widgets.add(Padding(
          padding: const EdgeInsets.symmetric(horizontal: 8.0),
          child: InkWell(
            borderRadius: const BorderRadius.all(Radius.circular(8)),
            onTap: () {},
            onDoubleTap: () {
              context.read<FSBloc>().add(FSLocationChangeRequest(
                  location: locationSplit.sublist(0, i + 1).join('/')
              ));
            },
            onSecondaryTap: () => showFolderDetails(locationSplit.sublist(0, i + 1).join('/')),
            child: Padding(
              padding: const EdgeInsets.all(8.0),
              child: Center(
                  child: Text(
                      locationSplit[i],
                    style: const TextStyle(
                      color: Color.fromARGB(255, 20, 30, 70),
                      fontWeight: FontWeight.bold,
                      fontSize: 16,
                    ),
                  )
              ),
            ),
          ),
        ));

        widgets.add(const Icon(
          Icons.navigate_next,
          size: 16,
        ));
      }
      widgets.removeLast();
      return widgets;
    }

    return Expanded(
      child: SizedBox(
        height: 40,
        child: ListView(
          // reverse: true,
          shrinkWrap: true,
          scrollDirection: Axis.horizontal,
          children: getFolderPath(),
        ),
      ),
    );
  }
}