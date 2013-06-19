# procedure DissectRecord(RecordDecoder decoder, FieldWriter writer, int repetitionLevel):
#     Add current repetitionLevel and definition level to writer
#
#     seenFields = {} // empty set of integers
#
#     while decoder has more field values
#         FieldWriter chWriter = child of writer for field read by decoder
#
#         int chRepetitionLevel = repetitionLevel
#
#         if set seenFields contains field ID of chWriter
#             chRepetitionLevel = tree depth of chWriter
#         else
#             Add field ID of chWriter to seenFields
#         end if
#
#         if chWriter corresponds to an atomic field
#             Write value of current field read by decoder using chWriter at chRepetitionLevel
#         else
#             DissectRecord(new RecordDecoder for nested record read by decoder, chWriter, chRepetitionLevel)
#         end if
#     end while
# end procedure


def disect_record(record, schema, columns, rlevel=0):
    for key, child_schema in schema.iteritems():
        path = child_schema.get('path')
        if child_schema['type'] == int:
            print key, path
            record_value = record.get(path)
            dlevel = len(path.split('.'))
            dlevel = dlevel if record_value else dlevel - 1
            column = (record_value, rlevel, dlevel)
            if record_value:
                if child_schema.get('repeated'):
                    for val in record_value:
                        column = (val, rlevel, dlevel)
                        columns[path].append(column)
                else:
                    columns[path].append(column)
            else:
                columns[path].append((None, rlevel, dlevel))


        if child_schema['type'] == 'group':
            disect_record(record, child_schema['children'], columns)

    return columns



if __name__ == '__main__':

    schema = {
        'id': {'required': True, 'type': int, 'path': 'id'},
        'links': {'required': False, 'type': 'group', 'children': {
            'forward': {'required': False, 'type': int, 'repeated': True, 'path': 'links.forward'},
            'backward': {'required': False, 'type': int, 'repeated': True, 'path': 'links.backward'}
        }}
    }

    columns = {
        'id': [],
        'links.forward': [],
        'links.backward': []
    }

    r1 = {
        'id': 10,
        'links.forward': [20, 40, 80]
    }

    r2 = {
        'id': 20,
        'links.backward': [10, 30],
        'links.forward': [80]
    }

    print disect_record(r1, schema, columns)
    print '`````'
    print disect_record(r2, schema, columns)
