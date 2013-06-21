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


class Writer():
    def __init__(self, parent=None, field_id='__base__', actually_present=False):
        self.parent = parent
        self.field_id = field_id
        self.actually_present = actually_present # Is this field optional or repeated and actually present

    def write(self, value, rlevel):
        print self.path
        print (value, rlevel, self.get_dlevel(value))

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

    def get_full_tree_depth(self):
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
        child_writer = Writer(parent=writer, field_id=field.name, actually_present=(field.mode != 'required' and value))
        child_rlevel = rlevel
        if child_writer.field_id in seen_fields:
            child_rlevel = child_writer.get_full_tree_depth()
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
