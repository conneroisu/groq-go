library ieee;

use ieee.STD_LOGIC_1164.all;

entity OnesComp is

    generic (
        N : integer := 8
        );
    port (
        Input  : in  std_logic_vector(N-1 downto 0);
        Output : out std_logic_vector(N-1 downto 0)
        );

end OnesComp;

architecture Behavioral of OnesComp is
begin

    Complement_Process : process(Input)
    begin
        for i in 0 to N-1 loop
            Output(i) <= not Input(i);
        end loop;
    end process;

end Behavioral;
