import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:fs_frontend/blocs/transfer/transfer_bloc.dart';
import 'package:fs_frontend/utilities/file_size_utilities.dart';

class TransferPane extends StatelessWidget {
  const TransferPane({super.key});

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      height: 300,
      width: 500,
      child: BlocBuilder<TransferBloc, TransferState>(
        builder: (context, state) {
          final transfers = context.read<TransferBloc>().state.transfers;
          if (transfers.isEmpty) {
            return const Center(
              child: Text('No transfers currently in progress'),
            );
          }
          return ListView.builder(
            itemCount: transfers.length,
            itemBuilder: (context, index) {
              return Card(
                child: ListTile(
                  leading: SizedBox(
                    height: 20,
                    width: 20,
                    child: StreamBuilder<double>(
                      builder: (_, __) {
                        if (transfers[index].hasCompleted) {
                          return const Icon(Icons.done_all);
                        }
                        if (transfers[index].hasFailed) {
                          return const Icon(Icons.error_outline,
                              color: Colors.redAccent);
                        }
                        if (transfers[index].progress == 0) {
                          return const CircularProgressIndicator(strokeWidth: 2.5);
                        }
                        return Stack(
                          alignment: Alignment.center,
                          children: [
                            CircularProgressIndicator(
                                value: transfers[index].progress,
                                strokeWidth: 2.5),
                            Text(
                                (transfers[index].progress * 100)
                                    .round()
                                    .toString(),
                                style: const TextStyle(fontSize: 8))
                          ],
                        );
                      },
                      stream: transfers[index].progressStream,
                    ),
                  ),
                  title: Text(transfers[index].fileName,
                      maxLines: 1, overflow: TextOverflow.ellipsis),
                  subtitle: StreamBuilder<double>(
                    builder: (_, __) {
                      if (transfers[index].hasCompleted) {
                        return const Text('Done');
                      }
                      if (transfers[index].hasFailed) {
                        return const Text('Failed',
                            style: TextStyle(color: Colors.redAccent));
                      }
                      final completedSize = getSizeString(
                          (transfers[index].progress * transfers[index].size)
                              .round());
                      final totalSize = getSizeString(transfers[index].size);
                      return Text('$completedSize of $totalSize');
                    },
                    stream: transfers[index].progressStream,
                  ),
                ),
              );
            },
          );
        },
      ),
    );
  }
}

void showTransferPane(BuildContext context) {
  showDialog(
    context: context,
    builder: (_) => BlocProvider.value(
      value: context.read<TransferBloc>(),
      child: AlertDialog(
        title: const Text('Transfers'),
        content: const Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TransferPane(),
          ],
        ),
        actions: context.read<TransferBloc>().state.transfers.isEmpty
            ? [
                TextButton(
                    onPressed: () {
                      Navigator.of(context).pop();
                    },
                    child: const Text('Close'))
              ]
            : [
                TextButton(
                    onPressed: () {
                      context.read<TransferBloc>().add(ClearCompleted());
                    },
                    child: const Text('Clear Completed')),
                TextButton(
                    onPressed: () {
                      Navigator.of(context).pop();
                    },
                    child: const Text('Close'))
              ],
      ),
    ),
  );
}
