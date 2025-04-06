package runtime

func InitGlobalEnvironment() {
	initNode()
	initAliasDeclarationEntry()
	initAliasDeclarationNode()
	initAnyTypeNode()
	initArrayListLiteralNode()
	initArrayTupleLiteralNode()
	initAsExpressionNode()
	initAsPatternNode()
	initAssignmentExpressionNode()
	initAttrDeclarationNode()
	initAttributeAccessNode()
	initAttributeParameterNode()
	initAwaitExpressionNode()
	initBigFloatLiteralNode()
	initBinArrayListLiteralNode()
	initBinArrayTupleLiteralNode()
	initBinHashSetLiteralNode()
	initBinaryExpressionNode()
	initBinaryPatternNode()
	initBinaryTypeNode()
	initBoolLiteralNode()
	initBreakExpressionNode()
	initCallNode()
	initCaseNode()
	initCatchNode()
	initCharLiteralNode()
	initClassDeclarationNode()
	initClosureTypeNode()
	initClosureLiteralNode()
	initConstantAsNode()
	initConstantDeclarationNode()
	initConstantLookupNode()
	initConstructorCallNode()
	initContinueExpressionNode()
	initDoExpressionNode()
	initMacroBoundaryNode()
	initQuoteExpressionNode()
	initUnquoteExpressionNode()
	initDoubleQuotedStringLiteralNode()
	initDoubleSplatExpressionNode()
	initEmptyStatementNode()
	initExpressionStatementNode()
	initExtendWhereBlockExpressionNode()
	initTrueLiteralNode()
	initFalseLiteralNode()
	initFloat32LiteralNode()
	initFloat64LiteralNode()
	initFloatLiteralNode()
	initForInExpressionNode()
	initFormalParameterNode()
	initGenericConstantNode()
	initGenericConstructorCallNode()
	initGenericMethodCallNode()
	initGenericReceiverlessMethodCallNode()
	initGenericTypeDefinitionNode()
	initGetterDeclarationNode()
	initGoExpressionNode()
	initHashMapLiteralNode()
	initHashRecordLiteralNode()
	initHashSetLiteralNode()
	initHexArrayListLiteralNode()
	initHexArrayTupleLiteralNode()
	initHexHashSetLiteralNode()
	initIfExpressionNode()
	initImplementExpressionNode()
	initImportStatementNode()
	initIncludeExpressionNode()
	initInitDefinitionNode()
	initInstanceOfTypeNode()
	initInstanceVariableDeclarationNode()
	initInstanceVariableNode()
	initInt8LiteralNode()
	initInt16LiteralNode()
	initInt32LiteralNode()
	initInt64LiteralNode()
	initIntLiteralNode()
	initInterfaceDeclarationNode()
	initInterpolatedRegexLiteralNode()
	initInterpolatedStringLiteralNode()
	initInterpolatedSymbolLiteralNode()
	initIntersectionTypeNode()
	initInvalidNode()
	initKeyValueExpressionNode()
	initKeyValuePatternNode()
	initLabeledExpressionNode()
	initListPatternNode()
	initLogicalExpressionNode()
	initLoopExpressionNode()
	initMapPatternNode()
	initMethodCallNode()
	initMethodDefinitionNode()
	initMethodLookupAsNode()
	initMethodLookupNode()
	initMethodParameterNode()
	initMethodSignatureDefinitionNode()
	initMixinDeclarationNode()
	initModifierForInNode()
	initModifierIfElseNode()
	initModifierNode()
	initModuleDeclarationNode()
	initMustExpressionNode()
	initNamedCallArgumentNode()
	initNeverTypeNode()
	initNewExpressionNode()
	initNilLiteralNode()
	initNilSafeSubscriptExpressionNode()
	initNilableTypeNode()
	initNotTypeNode()
	initNumericForExpressionNode()
	initObjectPatternNode()
	initParameterStatementNode()
	initPostfixExpressionNode()
	initPrivateConstantNode()
	initPrivateIdentifierNode()
	initProgramNode()
	initPublicConstantAsNode()
	initPublicConstantNode()
	initPublicIdentifierAsNode()
	initPublicIdentifierNode()
	initRangeLiteralNode()
	initRawCharLiteralNode()
	initRawStringLiteralNode()
	initReceiverlessMethodCallNode()
	initRecordPatternNode()
	initRegexInterpolationNode()
	initRegexLiteralContentSectionNode()
	initRestPatternNode()
	initReturnExpressionNode()
	initSelfLiteralNode()
	initSetPatternNode()
	initSetterDeclarationNode()
	initSignatureParameterNode()
	initSimpleSymbolLiteralNode()
	initSingletonBlockExpressionNode()
	initSingletonTypeNode()
	initSplatExpressionNode()
	initStringInspectInterpolationNode()
	initStringInterpolationNode()
	initStringLiteralContentSectionNode()
	initStructDeclarationNode()
	initSubscriptExpressionNode()
	initSwitchExpressionNode()
	initSymbolArrayListLiteralNode()
	initSymbolArrayTupleLiteralNode()
	initSymbolHashSetLiteralNode()
	initSymbolKeyValueExpressionNode()
	initSymbolKeyValuePatternNode()
	initThrowExpressionNode()
	initTryExpressionNode()
	initTuplePatternNode()
	initTypeDefinitionNode()
	initTypeExpressionNode()
	initTypeofExpressionNode()
	initUInt8LiteralNode()
	initUInt16LiteralNode()
	initUInt32LiteralNode()
	initUInt64LiteralNode()
	initUnaryExpressionNode()
	initUnaryTypeNode()
	initUndefinedLiteralNode()
	initUninterpolatedRegexLiteralNode()
	initUnionTypeNode()
	initUnlessExpressionNode()
	initUntilExpressionNode()
	initUsingAllEntryNode()
	initUsingEntryWithSubentriesNode()
	initUsingExpressionNode()
	initValueDeclarationNode()
	initValuePatternDeclarationNode()
	initVariableDeclarationNode()
	initVariablePatternDeclarationNode()
	initVariantTypeParameterNode()
	initVoidTypeNode()
	initWhileExpressionNode()
	initWordArrayListLiteralNode()
	initWordArrayTupleLiteralNode()
	initWordHashSetLiteralNode()
	initYieldExpressionNode()
}

func init() {
	InitGlobalEnvironment()
}
