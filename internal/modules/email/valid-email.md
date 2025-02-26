## Regras de negócios para a feature de validação de E-mail

# Table email_validations

- id
- email_id
- attempts
- is_consumed
- is_valid
- created
- expires

### Regras de Validação

[] - Deve ser possível validar apenas um e-mail por vez
[] - Deve lançar um erro caso já exista uma validação ativa para o mesmo e-mail
[] - O código de validação deve expirar após 30 minutos
[] - O usuário deve ter no máximo 3 tentativas de validação
[] - Após 3 tentativas inválidas, a validação deve ser marcada como consumida
[] - Não deve ser possível validar um código expirado
[] - Após validação bem-sucedida, o e-mail deve ser marcado como verificado
[] - O código de validação deve ter 6 dígitos numéricos
[] - Não deve ser possível criar nova validação se já existir uma validação não expirada
[] - O sistema deve enviar e-mail com o código de validação automaticamente

### Estados possíveis da validação

- Pendente: Aguardando validação do usuário
- Expirada: Código não foi validado dentro do tempo limite
- Consumida: Todas as tentativas foram utilizadas ou validação foi concluída
- Válida: E-mail foi validado com sucesso

### Fluxo de validação

1. Usuário solicita validação de e-mail
2. Sistema verifica se não há validação ativa
3. Sistema gera código de 6 dígitos
4. Sistema envia e-mail com código
5. Sistema cria registro na tabela com status pendente
6. Usuário insere código
7. Sistema valida código e atualiza status
