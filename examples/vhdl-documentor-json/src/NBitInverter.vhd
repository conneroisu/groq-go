library IEEE;

use IEEE.STD_LOGIC_1164.all;

entity NBitInverter is

    generic (
        N : integer := 4                -- N is the width of the input/output
        );
    port (
        Input  : in  std_logic_vector(N-1 downto 0);  -- Input vector
        Output : out std_logic_vector(N-1 downto 0)   -- Output vector
        );

end NBitInverter;

architecture Behavioral of NBitInverter is
begin

    Complement_Process : process(Input)
    begin
        for i in 0 to N-1 loop
            if Input(i) = '1' then
                Output(i) <= '0';
            else
                Output(i) <= '1';
            end if;
        end loop;
    end process;

end Behavioral;
