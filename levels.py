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
        self.rlevel = None
        self.dlevel = None

    def write(self, value):
        print self.path
        print (value, self.rlevel, self.dlevel)

    def get_tree_depth(self):
        depth = 0
        parent = self.parent
        while parent:
            depth += 1
            parent = parent.parent
        return depth

    @property
    def path(self):
        parts = [self.field_id]
        parent = self.parent
        while parent:
            parts.append(parent.field_id)
            parent = parent.parent
        ".".join(parts)

class Decoder():
    def __init__(self, schema, record):
        self.schema = schema
        self.record = record

    @property
    def fields(self):
        print self.schema
        return self.schema.fields

    def get_value(self, field_id):
        # Need to change this depending on whether value is atomic or not?
        # Or will we never have a decoder with an atomic value?
        return self.record.get(field_id)

# dlevel = how many fields optional or repeated fields are actually present
# rlevel = The level of the last repeated field

def findit(field_id, fields):
    matches = [f for f in fields if f.name == field_id]
    if len(matches) == 1:
        return matches[0]

def stripe_record(decoder, writer, rlevel=0):
    writer.rlevel = rlevel
    seen_fields = set()

    for field in decoder.fields:
        child_writer = Writer(parent=writer, field_id=field.name) # child of `writer` for field read by `decoder`
        child_rlevel = rlevel
        if child_writer.field_id in seen_fields:
            child_rlevel = child_writer.get_tree_depth()
        else:
            seen_fields.add(child_writer.field_id)

        if field.type != 'record':
            # How to we handle repeated atomic values?
            child_writer.write(decoder.get_value(field.name))
        else:
            child_schema = findit(field.name, decoder.schema.fields)
            print child_schema
            child_decoder = Decoder(child_schema, decoder.get_value(field.name))
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

    print stripe_record(Decoder(schema, r1), Writer())
    #r1 = {
        #'id': 10,
        #'links.forward': [20, 40, 80]
    #}

    #r2 = {
        #'id': 20,
        #'links.backward': [10, 30],
        #'links.forward': [80]
    #}

    #columns = {
        #'id': [],
        #'links.forward': [],
        #'links.backward': []
    #}

    #disect_record(r1, schema, columns)
    #print '`````'
    #disect_record(r2, schema, columns)

    #print columns['id']
    #print columns['links.forward']
    #print columns['links.backward']



