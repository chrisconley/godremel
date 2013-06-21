import re
from collections import namedtuple
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

class Writer():
    def __init__(self, parent=None, field_id='__base__'):
        self.parent = parent
        self.field_id = field_id

    def write(self, value, rlevel):
        print self.path
        print (value, rlevel, self.get_dlevel(value))

    def get_dlevel(self, value):
        depth = self.get_tree_depth()
        return depth + 1 if value else depth

    def get_tree_depth(self):
        depth = 0
        parent = self.parent
        while parent:
            parent = parent.parent
            depth += 1
        return depth - 1

    @property
    def path(self):
        parts = [self.field_id]
        parent = self.parent
        while parent:
            parts.append(parent.field_id)
            parent = parent.parent
        parts.reverse()
        return ".".join(parts)

class Decoder():
    def __init__(self, schema, record):
        self.schema = schema
        self.record = record

    @property
    def fields(self):
        return self.schema.fields

    def field_values(self):
        for field in self.fields:
            if field.mode == 'repeated' and self.get_value(field.name):
                for value in self.get_value(field.name):
                    yield field, value
            else:
                yield field, self.get_value(field.name)

    def get_value(self, field_id):
        # Need to change this depending on whether value is atomic or not?
        # Or will we never have a decoder with an atomic value?
        return self.record and self.record.get(field_id)

# dlevel = how many fields optional or repeated fields are actually present
# rlevel = The level of the last repeated field

def findit(field_id, fields):
    matches = [f for f in fields if f.name == field_id]
    if len(matches) == 1:
        return matches[0]

def stripe_record(decoder, writer, rlevel=0):
    # add current definition and repetition level to writer (something to do with version maybe)?
    seen_fields = set()

    for field, value in decoder.field_values(): # in field *values*?!
        child_writer = Writer(parent=writer, field_id=field.name) # child of `writer` for field read by `decoder`
        child_rlevel = rlevel
        if child_writer.field_id in seen_fields:
            child_rlevel = child_writer.get_tree_depth()
        else:
            seen_fields.add(child_writer.field_id)

        if field.type != 'record':
            child_writer.write(value, child_rlevel)
        else:
            child_schema = findit(field.name, decoder.schema.fields)
            child_decoder = Decoder(child_schema, value)
            stripe_record(child_decoder, child_writer, child_rlevel)


Field = namedtuple('Field', 'name type mode fields')

if __name__ == '__main__':

    schema = Field('__base__', 'record', 'required',
        [
            Field('id', int, 'required', None),
            Field('links', 'record', 'optional', [
                Field('forward', int, 'repeated', None),
                Field('backward', int, 'repeated', None),
            ]),
        ]
    )

    r1 = {
        'id': 10,
        'links' : {
            'forward': [20, 40, 60]
        }
    }

    r2 = {
        'id': 20,
        'links' : {
            'backward': [10, 30],
            'forward': [80]
        }
    }

    r3 = {
        'id': 30,
    }

    print stripe_record(Decoder(schema, r1), Writer())
    print '````'
    print stripe_record(Decoder(schema, r2), Writer())
    print '````'
    print stripe_record(Decoder(schema, r3), Writer())
