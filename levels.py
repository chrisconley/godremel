import re
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
    seen_fields = set()
    for key, child_schema in schema.iteritems():
        path = child_schema.get('path')
        child_rlevel = rlevel

        if child_schema['type'] == int:
            record_value = record.get(path)

            if re.search('\.', path):
                dlevel = len(path.split('.'))
                dlevel = dlevel if record_value else dlevel - 1
            else:
                dlevel = 0
            column = (record_value, child_rlevel, dlevel)
            if record_value:
                if child_schema.get('repeated'):
                    for val in record_value:
                        if path in seen_fields:
                            child_rlevel = rlevel +  1
                        else:
                            seen_fields.add(path)
                        column = (val, child_rlevel, dlevel)
                        columns[path].append(column)
                else:
                    columns[path].append(column)
            else:
                columns[path].append((None, child_rlevel, dlevel))


        if child_schema['type'] == 'group':
            disect_record(record, child_schema['children'], columns, child_rlevel)

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

    disect_record(r1, schema, columns)
    print '`````'
    disect_record(r2, schema, columns)

    print columns['id']
    print columns['links.forward']
    print columns['links.backward']



