import re
from collections import namedtuple, defaultdict
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

columns = defaultdict(list)


class Writer():
    def __init__(self, parent=None, field=None, value=None):
        self.parent = parent
        self.field = field
        self.field_id = field.name
        self.value = value

        # Is this field optional or repeated and actually present
        self.actually_present= (field.mode != 'required' and value)

        # Is this a repeated field
        self.repeated = field.mode == 'repeated'

    def write(self, value, rlevel):
        column = (value, rlevel, self.get_dlevel(value))
        columns[self.path].append(column)

    def get_dlevel(self, value):
        depth = self.get_tree_depth()
        return depth + 1 if self.actually_present else depth

    def get_tree_depth(self):
        depth = 0
        parent = self.parent
        while parent:
            if parent and parent.actually_present:
                depth += 1
            parent = parent.parent
        return depth

    def get_repeated_field_depth(self):
        depth = 1 # for self
        #depth = 1 if self.repeated else 0
        parent = self.parent
        while parent:
            if parent and parent.repeated:
                depth += 1
            parent = parent.parent
        return depth

    def get_full_tree_depth(self):
        #depth = 1 # for self
        depth = 0
        parent = self.parent
        while parent:
            parent = parent.parent and parent.mode == 'repeated'
            #if parent:
                #depth += 1
            depth += 1
        return depth

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

    for field, value in decoder.field_values():
        child_writer = Writer(parent=writer, field=field, value=value)
        child_rlevel = rlevel
        if child_writer.field_id in seen_fields:
            child_rlevel = child_writer.get_repeated_field_depth()
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
            Field('names', 'record', 'repeated', [
                Field('languages', 'record', 'repeated', [
                    Field('code', str, 'required', None),
                    Field('country', str, 'optional', None)
                ]),
                Field('url', str, 'optional', None)
            ])

        ]
    )

    r1 = {
        'id': 10,
        'links' : {
            'forward': [20, 40, 60]
        },
        'names': [
            {
                'languages': [
                    {'code': 'en-us', 'country': 'us'},
                    {'code': 'en'}
                ],
                'url': 'http://A'
            },
            {
                'url': 'http://B'
            },
            {
                'languages': [
                    {'code': 'en-gb', 'country': 'en'}
                ]
            }
        ]
    }

    r2 = {
        'id': 20,
        'links' : {
            'backward': [10, 30],
            'forward': [80]
        },
        'names': [
            {'url': 'http://C'}
        ]
    }

    r3 = {
        'id': 30,
    }

    stripe_record(Decoder(schema, r1), Writer(field=schema))
    stripe_record(Decoder(schema, r2), Writer(field=schema))
    for path, entries in columns.iteritems():
        print path
        for e in entries:
            print e
