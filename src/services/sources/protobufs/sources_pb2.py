# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: sources.proto
# Protobuf Python Version: 5.27.2
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    5,
    27,
    2,
    '',
    'sources.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\rsources.proto\x12\x07sources\"9\n\x10OrderInfoRequest\x12\x10\n\x08order_id\x18\x01 \x01(\t\x12\x13\n\x0b\x65xecutor_id\x18\x02 \x01(\t\"\xc3\x01\n\x11OrderInfoResponse\x12\x10\n\x08order_id\x18\x01 \x01(\t\x12\x19\n\x11\x66inal_coin_amount\x18\x02 \x01(\x02\x12\x32\n\x10price_components\x18\x03 \x01(\x0b\x32\x18.sources.PriceComponents\x12\x32\n\x10\x65xecutor_profile\x18\x04 \x01(\x0b\x32\x18.sources.ExecutorProfile\x12\x19\n\x11zone_display_name\x18\x05 \x01(\t\";\n\x0f\x45xecutorProfile\x12\n\n\x02id\x18\x01 \x01(\t\x12\x0c\n\x04tags\x18\x02 \x03(\t\x12\x0e\n\x06rating\x18\x03 \x01(\x02\"U\n\x0fPriceComponents\x12\x18\n\x10\x62\x61se_coin_amount\x18\x01 \x01(\x02\x12\x12\n\ncoin_coeff\x18\x02 \x01(\x02\x12\x14\n\x0c\x62onus_amount\x18\x03 \x01(\x02\x32Y\n\x10OrderInfoService\x12\x45\n\x0cGetOrderInfo\x12\x19.sources.OrderInfoRequest\x1a\x1a.sources.OrderInfoResponseb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'sources_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  DESCRIPTOR._loaded_options = None
  _globals['_ORDERINFOREQUEST']._serialized_start=26
  _globals['_ORDERINFOREQUEST']._serialized_end=83
  _globals['_ORDERINFORESPONSE']._serialized_start=86
  _globals['_ORDERINFORESPONSE']._serialized_end=281
  _globals['_EXECUTORPROFILE']._serialized_start=283
  _globals['_EXECUTORPROFILE']._serialized_end=342
  _globals['_PRICECOMPONENTS']._serialized_start=344
  _globals['_PRICECOMPONENTS']._serialized_end=429
  _globals['_ORDERINFOSERVICE']._serialized_start=431
  _globals['_ORDERINFOSERVICE']._serialized_end=520
# @@protoc_insertion_point(module_scope)